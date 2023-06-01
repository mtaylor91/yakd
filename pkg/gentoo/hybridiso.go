package gentoo

import (
	"context"
	"fmt"
	"io"
	"os"
	"path"
	"path/filepath"

	log "github.com/sirupsen/logrus"

	"github.com/mtaylor91/yakd/pkg/util"
	"github.com/mtaylor91/yakd/pkg/util/executor"
)

const grubCfg = `
serial --speed=9600 --unit=0 --word=8 --parity=no --stop=1
terminal_input console serial
terminal_output console serial

set default=0
set timeout=5

menuentry "YAKD" {
	linux /boot/vmlinuz rd.live.debug root=live:CDLABEL=YAKD console=tty0 console=ttyS0,9600n8
	initrd /boot/initramfs
}
`

type HybridISOBuilder struct {
	ISODir string
	Target string
}

func (g *HybridISOBuilder) BuildISO(ctx context.Context) error {
	// Build ISO
	if err := executor.RunCmd(ctx, "xorrisofs",
		"-iso-level", "3",
		"-full-iso9660-filenames",
		"-volid", "YAKD",
		"-eltorito-boot", "bios.img",
		"-no-emul-boot", "-boot-load-size", "4", "-boot-info-table",
		"-isohybrid-mbr",
		path.Join(g.ISODir, "isohdpfx.bin"),
		"--efi-boot", "efi.img",
		"-efi-boot-part",
		"--efi-boot-image",
		"--protective-msdos-label",
		"-output", g.Target,
		g.ISODir,
	); err != nil {
		return err
	}

	return nil
}

type HybridISOSourceBuilder struct {
	BinPkgsCache string
	FSDir        string
	ISODir       string
}

func (g *HybridISOSourceBuilder) BuildISOFS(
	ctx context.Context, chroot executor.Executor,
) error {
	// Bind BinPkgsCache to /var/cache/binpkgs
	if err := executor.Default.RunCmd(
		ctx, "mount", "--bind",
		g.BinPkgsCache,
		path.Join(g.FSDir, "var/cache/binpkgs"),
	); err != nil {
		return err
	}

	// Unmount /var/cache/binpkgs on exit
	defer func() {
		if err := executor.Default.RunCmd(
			ctx, "umount", path.Join(g.FSDir, "var/cache/binpkgs"),
		); err != nil {
			log.Warnf("Failed to unmount /var/cache/binpkgs: %s", err)
		}
	}()

	// Install grub
	if err := installPackages(ctx, chroot, "sys-boot/grub"); err != nil {
		return err
	}

	// Install syslinux
	if err := installPackages(ctx, chroot, "sys-boot/syslinux"); err != nil {
		return err
	}

	// Install mkfs.vfat
	if err := installPackages(ctx, chroot, "sys-fs/dosfstools"); err != nil {
		return err
	}

	// Install squashfs-tools
	if err := installPackages(ctx, chroot, "sys-fs/squashfs-tools"); err != nil {
		return err
	}

	// Create FS/boot/build/esp
	esp := path.Join(g.FSDir, "boot", "build", "esp")
	if err := os.MkdirAll(esp, 0755); err != nil {
		return err
	}

	// Build grub MBR image
	if err := g.isoBuildBIOS(ctx, chroot); err != nil {
		return err
	}

	// Build grub EFI image
	if err := g.isoBuildEFI(ctx, esp, chroot); err != nil {
		return err
	}

	// Copy isohdpfx.bin to ISO directory
	if err := util.CopyFile(
		path.Join(g.FSDir, "usr", "share", "syslinux", "isohdpfx.bin"),
		path.Join(g.ISODir, "isohdpfx.bin"),
	); err != nil {
		return err
	}

	log.Warning("TODO: cleanup filesystem")

	// Identify kernel version
	kernelModules, err := filepath.Glob(
		filepath.Join(g.FSDir, "lib", "modules", "*"))
	if err != nil {
		return err
	}
	if len(kernelModules) != 1 {
		return fmt.Errorf("expected 1 kernel, found %d", len(kernelModules))
	}
	kernelVersion := filepath.Base(kernelModules[0])

	// Build initramfs
	if err := chroot.RunCmd(
		ctx, "dracut", "--force", "--kver", kernelVersion,
		"--add", "dmsquash-live", "--add", "pollcdrom"); err != nil {
		return err
	}

	return nil
}

func (g *HybridISOSourceBuilder) BuildISOSources(ctx context.Context) error {
	// Create LiveOS directory
	liveOS := path.Join(g.ISODir, "LiveOS")
	if err := os.MkdirAll(liveOS, 0755); err != nil {
		return err
	}

	// Build root squashfs
	if err := executor.RunCmd(ctx, "mksquashfs", g.FSDir,
		path.Join(liveOS, "squashfs.img")); err != nil {
		return err
	}

	// Copy built bootloader images to target directory
	log.Infof("Copying GRUB images to %s", g.ISODir)

	if err := util.CopyFile(
		filepath.Join(g.FSDir, "boot", "build", "bios.img"),
		filepath.Join(g.ISODir, "bios.img"),
	); err != nil {
		return err
	}

	if err := util.CopyFile(
		filepath.Join(g.FSDir, "boot", "build", "efi.img"),
		filepath.Join(g.ISODir, "efi.img"),
	); err != nil {
		return err
	}

	// Clean up build directory
	log.Infof("Cleaning up after GRUB image builds")
	err := os.RemoveAll(filepath.Join(g.FSDir, "boot", "build"))
	if err != nil {
		return err
	}

	log.Infof("Copying kernel and initramfs to %s", g.ISODir)

	// Identify kernel
	kernels, err := filepath.Glob(filepath.Join(g.FSDir, "boot", "vmlinuz-*"))
	if err != nil {
		return err
	}
	if len(kernels) != 1 {
		return fmt.Errorf("expected 1 kernel, found %d", len(kernels))
	}
	kernel := kernels[0]

	// Identify initramfs
	initramfss, err := filepath.Glob(filepath.Join(g.FSDir, "boot", "initramfs-*"))
	if err != nil {
		return err
	}
	if len(initramfss) != 1 {
		return fmt.Errorf("expected 1 initramfs, found %d", len(initramfss))
	}
	initramfs := initramfss[0]

	// Create boot directory
	if err := os.MkdirAll(filepath.Join(g.ISODir, "boot"), 0755); err != nil {
		return err
	}

	// Copy kernel
	if err := util.CopyFile(
		kernel,
		filepath.Join(g.ISODir, "boot", "vmlinuz"),
	); err != nil {
		return err
	}

	// Copy initramfs
	if err := util.CopyFile(
		initramfs,
		filepath.Join(g.ISODir, "boot", "initramfs"),
	); err != nil {
		return err
	}

	return nil
}

func (g *HybridISOSourceBuilder) isoBuildBIOS(
	ctx context.Context, chroot executor.Executor,
) error {
	log.Info("Building grub BIOS image")

	// Build grub BIOS image
	if err := chroot.RunCmd(ctx, "grub-mkimage", "-O", "i386-pc",
		"-o", path.Join("/boot", "build", "core.img"),
		"-p", "/boot/grub",
		"biosdisk", "fat", "iso9660", "part_gpt", "part_msdos", "normal", "boot",
		"linux", "configfile", "loopback", "chain", "ls", "search",
		"search_label", "search_fs_uuid", "search_fs_file", "gfxterm",
		"gfxterm_background", "gfxterm_menu", "test", "all_video", "loadenv",
		"exfat", "ext2", "ntfs", "btrfs", "hfsplus", "udf",
	); err != nil {
		return err
	}

	// Open <target>/usr/lib/grub/i386-pc/cdboot.img
	cdbootImg, err := os.Open(
		path.Join(g.FSDir, "usr", "lib", "grub", "i386-pc", "cdboot.img"))
	if err != nil {
		return err
	}

	defer cdbootImg.Close()

	// Open <target>/boot/build/core.img
	coreImg, err := os.Open(
		path.Join(g.FSDir, "boot", "build", "core.img"))
	if err != nil {
		return err
	}

	defer coreImg.Close()

	// Create /boot/build/bios.img
	biosImg, err := os.OpenFile(
		path.Join(g.FSDir, "boot", "build", "bios.img"),
		os.O_CREATE|os.O_RDWR, 0644)
	if err != nil {
		return err
	}

	defer biosImg.Close()

	// Concatenate cdboot.img and core.img
	if _, err := io.Copy(biosImg, cdbootImg); err != nil {
		return err
	}

	if _, err := io.Copy(biosImg, coreImg); err != nil {
		return err
	}

	// Create grub directory
	err = os.MkdirAll(filepath.Join(g.ISODir, "boot", "grub"), 0755)
	if err != nil {
		return err
	}

	// Write grub.cfg
	log.Infof("Writing GRUB configuration to %s/boot/grub/grub.cfg", g.ISODir)
	err = util.WriteFile(
		filepath.Join(g.ISODir, "boot", "grub", "grub.cfg"), grubCfg)
	if err != nil {
		return err
	}

	return nil
}

func (g *HybridISOSourceBuilder) isoBuildEFI(
	ctx context.Context, esp string, chroot executor.Executor,
) error {
	log.Info("Building grub EFI image")

	img := util.NewRawImage(
		path.Join(g.FSDir, "boot", "build", "efi.img"), 10, false)
	if err := img.Alloc(ctx); err != nil {
		return err
	}

	loop, err := img.Attach(ctx)
	if err != nil {
		return err
	}

	defer loop.Detach()

	// Create filesystem
	err = chroot.RunCmd(ctx, "mkfs.vfat", "-F", "16", loop.DevicePath)
	if err != nil {
		return err
	}

	// Mount filesystem
	if err := util.Mount(ctx, loop.DevicePath, esp); err != nil {
		return err
	}

	defer util.Unmount(ctx, esp)

	// Create /boot/build/esp/EFI/BOOT
	if err := os.MkdirAll(path.Join(esp, "EFI", "BOOT"), 0755); err != nil {
		return err
	}

	// Build grub EFI image
	if err := chroot.RunCmd(ctx, "grub-mkimage",
		"-O", "x86_64-efi",
		"-o", path.Join("/boot/build/esp/EFI/BOOT/BOOTX64.EFI"),
		"-p", "/boot/grub",
		"fat", "iso9660", "part_gpt", "part_msdos", "normal", "boot",
		"linux", "configfile", "loopback", "chain", "efifwsetup", "efi_gop",
		"efi_uga", "ls", "search", "search_label", "search_fs_uuid",
		"search_fs_file", "gfxterm", "gfxterm_background", "gfxterm_menu",
		"test", "all_video", "loadenv", "exfat", "ext2", "ntfs", "btrfs",
		"hfsplus", "udf",
	); err != nil {
		return err
	}

	return nil
}
