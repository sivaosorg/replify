package test

import (
	"os"
	"runtime"
	"strings"
	"testing"
	"time"

	"github.com/sivaosorg/replify/pkg/sysx"
)

// ///////////////////////////
// Section: OS detection tests
// ///////////////////////////

func TestSysx_OSDetection(t *testing.T) {
	t.Parallel()
	isLinux := sysx.IsLinux()
	isDarwin := sysx.IsDarwin()
	isWindows := sysx.IsWindows()

	count := 0
	if isLinux {
		count++
	}
	if isDarwin {
		count++
	}
	if isWindows {
		count++
	}
	if count != 1 {
		t.Errorf("expected exactly one OS flag to be true; got IsLinux=%v IsDarwin=%v IsWindows=%v",
			isLinux, isDarwin, isWindows)
	}
}

func TestSysx_OSName(t *testing.T) {
	t.Parallel()
	name := sysx.OSName()
	if name == "" {
		t.Error("OSName() returned empty string")
	}
}

func TestSysx_Arch(t *testing.T) {
	t.Parallel()
	arch := sysx.Arch()
	if arch == "" {
		t.Error("Arch() returned empty string")
	}
}

// ///////////////////////////
// Section: Runtime tests
// ///////////////////////////

func TestSysx_PID(t *testing.T) {
	t.Parallel()
	pid := sysx.PID()
	if pid <= 0 {
		t.Errorf("PID() = %d; want positive integer", pid)
	}
}

func TestSysx_UID(t *testing.T) {
	t.Parallel()
	uid := sysx.UID()
	if uid < 0 && runtime.GOOS != "windows" {
		t.Errorf("UID() = %d; want non-negative integer on non-Windows", uid)
	}
}

func TestSysx_GID(t *testing.T) {
	t.Parallel()
	gid := sysx.GID()
	if gid < 0 && runtime.GOOS != "windows" {
		t.Errorf("GID() = %d; want non-negative integer on non-Windows", gid)
	}
}

func TestSysx_Hostname(t *testing.T) {
	t.Parallel()
	host, err := sysx.Hostname()
	if err != nil {
		t.Fatalf("Hostname() error = %v", err)
	}
	if host == "" {
		t.Error("Hostname() returned empty string")
	}
}

func TestSysx_GoVersion(t *testing.T) {
	t.Parallel()
	v := sysx.GoVersion()
	if !strings.HasPrefix(v, "go") {
		t.Errorf("GoVersion() = %q; want string starting with 'go'", v)
	}
}

func TestSysx_NumCPU(t *testing.T) {
	t.Parallel()
	n := sysx.NumCPU()
	if n < 1 {
		t.Errorf("NumCPU() = %d; want >= 1", n)
	}
}

// ///////////////////////////
// Section: Environment tests
// ///////////////////////////

func TestSysx_GetEnv(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name     string
		key      string
		setup    func()
		teardown func()
		fallback string
		want     string
	}{
		{
			name:     "unset variable returns fallback",
			key:      "SYSX_TEST_UNSET_XYZ",
			setup:    func() { os.Unsetenv("SYSX_TEST_UNSET_XYZ") },
			teardown: func() {},
			fallback: "default",
			want:     "default",
		},
		{
			name: "set variable returns value",
			key:  "SYSX_TEST_SET_XYZ",
			setup: func() {
				os.Setenv("SYSX_TEST_SET_XYZ", "hello")
			},
			teardown: func() { os.Unsetenv("SYSX_TEST_SET_XYZ") },
			fallback: "default",
			want:     "hello",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()
			defer tt.teardown()
			got := sysx.GetEnv(tt.key, tt.fallback)
			if got != tt.want {
				t.Errorf("GetEnv(%q, %q) = %q; want %q", tt.key, tt.fallback, got, tt.want)
			}
		})
	}
}

func TestSysx_HasEnv(t *testing.T) {
	t.Parallel()
	const key = "SYSX_HAS_ENV_TEST"
	os.Unsetenv(key)
	if sysx.HasEnv(key) {
		t.Errorf("HasEnv(%q) = true; want false when not set", key)
	}
	os.Setenv(key, "yes")
	defer os.Unsetenv(key)
	if !sysx.HasEnv(key) {
		t.Errorf("HasEnv(%q) = false; want true when set", key)
	}
}

func TestSysx_GetEnvInt(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name     string
		value    string
		set      bool
		fallback int
		want     int
	}{
		{"unset returns fallback", "", false, 42, 42},
		{"valid int", "10", true, 42, 10},
		{"invalid value returns fallback", "abc", true, 42, 42},
	}
	const key = "SYSX_INT_TEST"
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.set {
				os.Setenv(key, tt.value)
				defer os.Unsetenv(key)
			} else {
				os.Unsetenv(key)
			}
			got := sysx.GetEnvInt(key, tt.fallback)
			if got != tt.want {
				t.Errorf("GetEnvInt(%q, %d) = %d; want %d", key, tt.fallback, got, tt.want)
			}
		})
	}
}

func TestSysx_GetEnvBool(t *testing.T) {
	t.Parallel()
	tests := []struct {
		value    string
		fallback bool
		want     bool
	}{
		{"1", false, true},
		{"true", false, true},
		{"yes", false, true},
		{"on", false, true},
		{"TRUE", false, true},
		{"YES", false, true},
		{"0", true, false},
		{"false", true, false},
		{"no", true, false},
		{"off", true, false},
		{"garbage", true, true},
	}
	const key = "SYSX_BOOL_TEST"
	for _, tt := range tests {
		t.Run(tt.value+"_fallback_"+func() string {
			if tt.fallback {
				return "true"
			}
			return "false"
		}(), func(t *testing.T) {
			os.Setenv(key, tt.value)
			defer os.Unsetenv(key)
			got := sysx.GetEnvBool(key, tt.fallback)
			if got != tt.want {
				t.Errorf("GetEnvBool(%q, %v) with value %q = %v; want %v",
					key, tt.fallback, tt.value, got, tt.want)
			}
		})
	}
}

func TestSysx_GetEnvSlice(t *testing.T) {
	t.Parallel()
	const key = "SYSX_SLICE_TEST"
	os.Setenv(key, "a,b,c")
	defer os.Unsetenv(key)

	got := sysx.GetEnvSlice(key, ",")
	if len(got) != 3 {
		t.Fatalf("GetEnvSlice len = %d; want 3", len(got))
	}
	if got[0] != "a" || got[1] != "b" || got[2] != "c" {
		t.Errorf("GetEnvSlice = %v; want [a b c]", got)
	}

	os.Unsetenv(key)
	if sysx.GetEnvSlice(key, ",") != nil {
		t.Error("GetEnvSlice on unset key should return nil")
	}
}

func TestSysx_EnvMap(t *testing.T) {
	t.Parallel()
	m := sysx.EnvMap()
	if m == nil {
		t.Error("EnvMap() returned nil")
	}
}

// ///////////////////////////
// Section: File system tests
// ///////////////////////////

func TestSysx_FileExists(t *testing.T) {
	t.Parallel()
	f, err := os.CreateTemp("", "sysx_test_*")
	if err != nil {
		t.Fatal(err)
	}
	path := f.Name()
	f.Close()
	defer os.Remove(path)

	if !sysx.FileExists(path) {
		t.Errorf("FileExists(%q) = false; want true", path)
	}
	if sysx.FileExists(path + "_nonexistent") {
		t.Errorf("FileExists(%q) = true; want false", path+"_nonexistent")
	}
}

func TestSysx_DirExists(t *testing.T) {
	t.Parallel()
	dir, err := os.MkdirTemp("", "sysx_dir_*")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(dir)

	if !sysx.DirExists(dir) {
		t.Errorf("DirExists(%q) = false; want true", dir)
	}
	if sysx.DirExists(dir + "_nonexistent") {
		t.Errorf("DirExists(%q) = true; want false", dir+"_nonexistent")
	}
}

func TestSysx_IsFile_IsDir(t *testing.T) {
	t.Parallel()
	f, err := os.CreateTemp("", "sysx_isfile_*")
	if err != nil {
		t.Fatal(err)
	}
	fpath := f.Name()
	f.Close()
	defer os.Remove(fpath)

	dir, err := os.MkdirTemp("", "sysx_isdir_*")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(dir)

	if !sysx.IsFile(fpath) {
		t.Errorf("IsFile(%q) = false; want true", fpath)
	}
	if sysx.IsDir(fpath) {
		t.Errorf("IsDir(%q) = true for a file; want false", fpath)
	}
	if sysx.IsFile(dir) {
		t.Errorf("IsFile(%q) = true for a dir; want false", dir)
	}
	if !sysx.IsDir(dir) {
		t.Errorf("IsDir(%q) = false; want true", dir)
	}
}

func TestSysx_FileSize(t *testing.T) {
	t.Parallel()
	content := []byte("hello sysx")
	f, err := os.CreateTemp("", "sysx_size_*")
	if err != nil {
		t.Fatal(err)
	}
	path := f.Name()
	f.Write(content)
	f.Close()
	defer os.Remove(path)

	size, err := sysx.FileSize(path)
	if err != nil {
		t.Fatalf("FileSize error = %v", err)
	}
	if size != int64(len(content)) {
		t.Errorf("FileSize = %d; want %d", size, len(content))
	}

	_, err = sysx.FileSize(path + "_nonexistent")
	if err == nil {
		t.Error("FileSize on non-existent path should return error")
	}
}

func TestSysx_TempDir(t *testing.T) {
	t.Parallel()
	tmp := sysx.TempDir()
	if tmp == "" {
		t.Error("TempDir() returned empty string")
	}
}

func TestSysx_HomeDir(t *testing.T) {
	t.Parallel()
	home, err := sysx.HomeDir()
	if err != nil {
		t.Fatalf("HomeDir() error = %v", err)
	}
	if home == "" {
		t.Error("HomeDir() returned empty string")
	}
}

func TestSysx_WorkingDir(t *testing.T) {
	t.Parallel()
	wd, err := sysx.WorkingDir()
	if err != nil {
		t.Fatalf("WorkingDir() error = %v", err)
	}
	if wd == "" {
		t.Error("WorkingDir() returned empty string")
	}
}

// ///////////////////////////
// Section: Process tests
// ///////////////////////////

func TestSysx_ProcessExists(t *testing.T) {
	t.Parallel()
	if !sysx.ProcessExists(os.Getpid()) {
		t.Errorf("ProcessExists(os.Getpid()) = false; want true")
	}
	if sysx.ProcessExists(-1) {
		t.Error("ProcessExists(-1) = true; want false")
	}
}

func TestSysx_CurrentProcessName(t *testing.T) {
	t.Parallel()
	name := sysx.CurrentProcessName()
	if name == "" {
		t.Error("CurrentProcessName() returned empty string")
	}
}

// ///////////////////////////
// Section: Command execution tests
// ///////////////////////////

func TestSysx_ExecCommand(t *testing.T) {
	t.Parallel()
	var cmd string
	var args []string
	if runtime.GOOS == "windows" {
		cmd = "cmd"
		args = []string{"/c", "echo", "hello"}
	} else {
		cmd = "echo"
		args = []string{"hello"}
	}
	if err := sysx.ExecCommand(cmd, args...); err != nil {
		t.Errorf("ExecCommand(%q) error = %v", cmd, err)
	}
}

func TestSysx_ExecOutput(t *testing.T) {
	t.Parallel()
	var cmd string
	var args []string
	if runtime.GOOS == "windows" {
		cmd = "cmd"
		args = []string{"/c", "echo", "hello"}
	} else {
		cmd = "echo"
		args = []string{"hello"}
	}
	out, err := sysx.ExecOutput(cmd, args...)
	if err != nil {
		t.Fatalf("ExecOutput error = %v", err)
	}
	if !strings.Contains(out, "hello") {
		t.Errorf("ExecOutput = %q; want output containing 'hello'", out)
	}
}

func TestSysx_ExecCommandWithTimeout(t *testing.T) {
	t.Parallel()
	var cmd string
	var args []string
	if runtime.GOOS == "windows" {
		cmd = "cmd"
		args = []string{"/c", "echo", "timeout_test"}
	} else {
		cmd = "echo"
		args = []string{"timeout_test"}
	}

	// Should complete well within 5 seconds.
	if err := sysx.ExecCommandWithTimeout(5*time.Second, cmd, args...); err != nil {
		t.Errorf("ExecCommandWithTimeout error = %v", err)
	}

	// A command that takes longer than the timeout should fail.
	var sleepCmd string
	var sleepArgs []string
	if runtime.GOOS == "windows" {
		sleepCmd = "cmd"
		sleepArgs = []string{"/c", "timeout", "10"}
	} else {
		sleepCmd = "sleep"
		sleepArgs = []string{"10"}
	}
	err := sysx.ExecCommandWithTimeout(100*time.Millisecond, sleepCmd, sleepArgs...)
	if err == nil {
		t.Error("ExecCommandWithTimeout: expected error on timeout, got nil")
	}
}

// ///////////////////////////
// Section: Entry-point tests
// ///////////////////////////

func TestSysx_SystemInfo(t *testing.T) {
	t.Parallel()
	info := sysx.SystemInfo()
	if info == nil {
		t.Fatal("SystemInfo() returned nil")
	}
	required := []string{"os", "arch", "hostname", "pid", "go_version", "executable", "num_cpu"}
	for _, key := range required {
		if _, ok := info[key]; !ok {
			t.Errorf("SystemInfo() missing key %q", key)
		}
	}
	if info["os"] == "" {
		t.Error("SystemInfo()[\"os\"] is empty")
	}
	if info["arch"] == "" {
		t.Error("SystemInfo()[\"arch\"] is empty")
	}
	if info["go_version"] == "" {
		t.Error("SystemInfo()[\"go_version\"] is empty")
	}
}
