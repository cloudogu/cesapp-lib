package core

import (
	"maps"
	"slices"
)

// Capability represent POSIX capabilities type.
//
// See docs at https://manned.org/capabilities.7
type Capability string

// DefaultCapabilities in the Cloudogu Ecosystem, regardless of which container runtime it runs in.
// This way the Ecosystem can guarantee a consistent behavior.
var DefaultCapabilities = []Capability{
	Chown, DacOverride, Fsetid, Fowner, Setgid, Setuid, Setpcap, NetBindService, Kill,
}

// AllCapabilities includes all capabilities that could possibly be set.
// The special value ALL is not included.
var AllCapabilities = []Capability{
	AuditControl, AuditRead, AuditWrite, BlockSuspend, Bpf, CheckpointRestore, Chown,
	DacOverride, Fowner, Fsetid, IpcLock, IpcOwner, Kill, Lease, LinuxImmutable, MacAdmin,
	MacOverride, Mknod, NetAdmin, NetBindService, NetBroadcast, NetRaw, Perfmon, Setfcap,
	Setgid, Setpcap, Setuid, SysAdmin, SysBoot, SysChroot, SysModule, SysNice, SysPAcct,
	SysPTrace, SysResource, SysTime, SysTtyCONFIG, Syslog, WakeAlarm,
}

// These capabilities' documentation contain abstracts of their respective manpage documentation and may refer to
// other man pages references f. e. as epoll(7)
const (
	// All is a special capability which can be used to add or drop all capabilities listed below.
	All = "ALL"
	// AuditControl enables and disables kernel auditing; changes auditing filter rules.
	// retrieves auditing status and filtering rules.
	AuditControl = "AUDIT_CONTROL"
	// AuditRead allows reading the audit log via a multicast netlink socket.
	AuditRead = "AUDIT_READ"
	// AuditWrite write records to kernel auditing log.
	AuditWrite = "AUDIT_WRITE"
	// BlockSuspend employs features that can block system suspend (epoll(7) EPOLLWAKEUP, /proc/sys/wake_lock).
	BlockSuspend = "BLOCK_SUSPEND"
	// Bpf employs privileged BPF operations and separates out BPF functionality from the overloaded SYS_ADMIN capability.
	Bpf = "BPF"
	// CheckpointRestore allows facilitating checkpoint/restore for non-root users.
	CheckpointRestore = "CHECKPOINT_RESTORE"
	// Chown makes arbitrary changes to file UIDs and GIDs (see chown(2)).
	Chown = "CHOWN"
	// DacOverride allows bypassing file read, write, and execute permission checks to implement discretionary access control.
	DacOverride = "DAC_OVERRIDE"
	// Fowner allows to bypass permission checks on operations that require file ownership (i. e. read, write, or execute).
	Fowner = "FOWNER"
	// Fsetid allows to avoid resetting a file's Setuid flag on modifying.
	Fsetid = "FSETID"
	// IpcLock allows a process to lock memory to swapping the content to disk.
	IpcLock = "IPC_LOCK"
	// IpcOwner bypasses permission checks for operations on System V IPC objects.
	IpcOwner = "IPC_OWNER"
	// Kill bypasses permission checks for sending signals (see kill(2)). This includes use of the ioctl(2) KDSIGACCEPT operation.
	Kill = "KILL"
	// Lease establishes leases on arbitrary files (see fcntl(2)).
	Lease = "LEASE"
	// LinuxImmutable Set the FS_APPEND_FL and FS_IMMUTABLE_FL inode flags (see FS_IOC_SETFLAGS(2const)).
	LinuxImmutable = "LINUX_IMMUTABLE"
	// MacAdmin allows MAC configuration or state changes. Implemented for the Smack Linux Security Module (LSM).
	MacAdmin = "MAC_ADMIN"
	// MacOverride Override Mandatory Access Control (MAC).  Implemented for the Smack LSM.
	MacOverride = "MAC_OVERRIDE"
	// Mknod creates special files using mknod(2).
	Mknod = "MKNOD"
	// NetAdmin allows processes to perform various network-related operations.
	NetAdmin = "NET_ADMIN"
	// NetBindService binds a socket to Internet domain privileged ports (port numbers less than 1024).
	NetBindService = "NET_BIND_SERVICE"
	// NetBroadcast (Unused) makes socket broadcasts, and listen to multicasts.
	NetBroadcast = "NET_BROADCAST"
	// NetRaw uses RAW and PACKET sockets; binds to any address for transparent proxying.
	NetRaw = "NET_RAW"
	// Perfmon employs various performance-monitoring mechanisms various BPF operations that have performance implications.
	Perfmon = "PERFMON"
	// Setfcap allows to set arbitrary capabilities on a file. Since Linux 5.12, this capability is also needed to map
	// user ID 0 in a new user namespace; see user_namespaces(7) for details.
	Setfcap = "SETFCAP"
	// Setgid makes arbitrary manipulations of process GIDs and supplementary GID list;
	//    - forge GID when passing socket credentials via UNIX domain sockets;
	//    - write a group ID mapping in a user namespace (see user_namespaces(7)).
	Setgid = "SETGID"
	// Setpcap allows to add or drop any capability from the calling thread's bounding set.
	Setpcap = "SETPCAP"
	// Setuid allows to make arbitrary manipulations of process UIDs (setuid(2), setreuid(2), setresuid(2), setfsuid(2))
	Setuid = "SETUID"
	//	SysAdmin allows performing a range of system administration operations. This is like super dangerous and should not be granted easily
	SysAdmin = "SYS_ADMIN"
	//SysBoot allows to use reboot(2) and kexec_load(2).
	SysBoot = "SYS_BOOT"
	// SysChroot allows to use chroot(2) or change mount linux namespaces using setns(2).
	SysChroot = "SYS_CHROOT"
	// SysModule loads and unloads kernel modules (see init_module(2) and delete_module(2))
	SysModule = "SYS_MODULE"
	// SysNice lowers the process nice value (nice(2), setpriority(2)) and change the nice value for arbitrary processes
	SysNice = "SYS_NICE"
	//SysPAcct switches process accounting with acct(2) on or off.
	SysPAcct = "SYS_PACCT"
	//SysPTrace allows to trace arbitrary processes using ptrace(2), or read/write memory of arbitrary processes.
	SysPTrace = "SYS_PTRACE"
	// SysResource allows a process to modify the resource limits specified in variety of uses cases.
	SysResource = "SYS_RESOURCE"
	// SysTime sets system clock (settimeofday(2), stime(2), adjtimex(2)); set real-time (hardware) clock
	SysTime = "SYS_TIME"
	// SysTtyCONFIG uses vhangup(2); employ various privileged ioctl(2) operations on virtual terminals.
	SysTtyCONFIG = "SYS_TTY_CONFIG"
	// Syslog performs privileged syslog(2) operations.
	Syslog = "SYSLOG"
	// WakeAlarm triggers something that will wake up the system.
	WakeAlarm = "WAKE_ALARM"
)

// Capabilities represent POSIX capabilities that can be added to or removed from a dogu.
//
// The fields Add and Drop will modify the default capabilities as provided by DefaultCapabilities. Add will append
// further capabilities while Drop will remove capabilities. The capability All can be used to add or remove all
// available capabilities.
//
// See DefaultCapabilities for the standard set being used in the Cloudogu Ecosystem.
//
// This example will result in the following capability list: DacOverride, Fsetid, Fowner, Setgid, Setuid, Setpcap, NetBindService, Kill, Syslog
//
//	"Capabilities": {
//	   "Drop": "Chown"
//	   "Add": "Syslog"
//	}
//
// This example will result in the following capability list: NetBindService
//
//	"Capabilities": {
//	   "Drop": ["All"],
//	   "Add": ["NetBindService", "Kill"]
//	}
type Capabilities struct {
	// Add contains the capabilities that should be allowed to be used in a container. This list is optional.
	Add []Capability `json:"Add,omitempty"`
	// Drop contains the capabilities that should be blocked from being used in a container. This list is optional.
	Drop []Capability `json:"Drop,omitempty"`
}

// Security defines security policies for the dogu. These fields can be used to reduce a dogu's attack surface.
//
// Example:
//
//	"Security": {
//	  "Capabilities": {
//	     "Drop": ["All"],
//	     "Add": ["NetBindService", "Kill"]
//	   },
//	  "RunAsNonRoot": true,
//	  "ReadOnlyRootFileSystem": true
//	}
type Security struct {
	// Capabilities sets the allowed and dropped capabilities for the dogu. The dogu should not use more than the
	// configured capabilities here, otherwise failure may occur at start-up or at run-time. This list is optional.
	Capabilities Capabilities `json:"Capabilities,omitempty"`
	// RunAsNonRoot indicates that the container must run as a non-root user. The dogu must support running as non-root
	// user otherwise the dogu start may fail. This flag is optional and defaults to false.
	RunAsNonRoot bool
	// ReadOnlyRootFileSystem mounts the container's root filesystem as read-only. The dogu must support accessing the
	// root file system by only reading otherwise the dogu start may fail. This flag is optional and defaults to false.
	ReadOnlyRootFileSystem bool
}

// CalcEffectiveCapabilities returns the actual capabilities after dropping and then adding the given capabilities
// to the given default capabilities.
// It can also handle the All meta-capability, so adding or dropping all capabilities can be done
// without listing every single capability directly.
func CalcEffectiveCapabilities(defaultCaps, capsToDrop, capsToAdd []Capability) []Capability {
	effectiveCaps := make(map[Capability]int)

	for _, defaultCap := range defaultCaps {
		// note this works since go 1.22 because iteration variables now can be used as unshared variable
		effectiveCaps[defaultCap] = 0 // we only use the map to check for keys, values don't matter
	}

	for _, dropCap := range capsToDrop {
		if dropCap == All {
			effectiveCaps = make(map[Capability]int)
			break
		}
		delete(effectiveCaps, dropCap)
	}

	for _, addCap := range capsToAdd {
		if addCap == All {
			// do a fast exit here because alternatives of slice-to-map conversion would be cumbersome
			return slices.Clone(AllCapabilities)
		}
		effectiveCaps[addCap] = 0 // we only use the map to check for keys, values don't matter
	}

	actualCaps := maps.Keys(effectiveCaps)
	return slices.Collect(actualCaps)
}
