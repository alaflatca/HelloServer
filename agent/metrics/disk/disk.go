package disk

import (
	"helloServer/agent/measure"
	"syscall"

	"github.com/pkg/errors"
)

type metric struct {
	defaultPath string
}

// metric 여러 개 일 경우 config 받아서 처리

func New() *metric {
	return &metric{defaultPath: "/"}
}

func (mt *metric) Process(ms *measure.Measure) error {
	fs := syscall.Statfs_t{}
	err := syscall.Statfs(mt.defaultPath, &fs)
	if err != nil {
		return errors.Wrap(err, "Failed to syscall Statfs")
	}
	diskAll := float64(fs.Blocks * uint64(fs.Bsize))
	diskAvail := float64(fs.Bavail * uint64(fs.Bsize))
	ms.Disk.All = diskAll / measure.GB
	ms.Disk.Avail = diskAvail / measure.GB
	ms.Disk.Used = (diskAll - diskAvail) / measure.GB
	ms.Disk.Usage = ((diskAll - diskAvail) / diskAll) * 100
	// /dev/sda

	return nil
}

func (mt *metric) Once(ms *measure.Measure) error {
	return nil
}

// statfs () 시스템 호출 은 마운트된 파일 시스템에 대한 정보를 반환합니다.
//  path는 마운트된 파일 시스템 내의 모든 파일의 경로 이름입니다.
// buf는 대략 다음과 같이 정의된 statfs 구조 에 대한 포인터입니다 .

// struct statfs {
// __fsword_t f_type;    /* Type of filesystem (see below) */
// __fsword_t f_bsize;   /* Optimal transfer block size */
// fsblkcnt_t f_blocks;  /* Total data blocks in filesystem */
// fsblkcnt_t f_bfree;   /* Free blocks in filesystem */
// fsblkcnt_t f_bavail;  /* Free blocks available to unprivileged user */
// fsfilcnt_t f_files;   /* Total inodes in filesystem */
// fsfilcnt_t f_ffree;   /* Free inodes in filesystem */
// fsid_t     f_fsid;    /* Filesystem ID */
// __fsword_t f_namelen; /* Maximum length of filenames */
// __fsword_t f_frsize;  /* Fragment size (since Linux 2.6) */
// __fsword_t f_flags;   /* Mount flags of filesystem (since Linux 2.6.36) */
// __fsword_t f_spare[xxx]; /* Padding bytes reserved for future use */
// };
