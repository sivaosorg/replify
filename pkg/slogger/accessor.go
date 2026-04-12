package slogger

import (
	"io"
	"time"
)

// ///////////////////////////////////////////////////////////////////////////
// Options accessors
// ///////////////////////////////////////////////////////////////////////////

// GetLevel returns the minimum log level.
//
// Returns:
//
// the Level configured in Options.
func (o *Options) GetLevel() Level {
	if o == nil {
		return InfoLevel
	}
	return o.level
}

// SetLevel sets the minimum log level.
//
// Parameters:
//   - `level`: the minimum log level
func (o *Options) SetLevel(level Level) {
	if o == nil {
		return
	}
	o.level = level
}

// GetFormatter returns the formatter.
//
// Returns:
//
// the Formatter configured in Options.
func (o *Options) GetFormatter() Formatter {
	if o == nil {
		return nil
	}
	return o.formatter
}

// SetFormatter sets the formatter.
//
// Parameters:
//   - `formatter`: the formatter to use
func (o *Options) SetFormatter(formatter Formatter) {
	if o == nil {
		return
	}
	o.formatter = formatter
}

// GetOutput returns the output writer.
//
// Returns:
//
// the io.Writer configured in Options.
func (o *Options) GetOutput() io.Writer {
	if o == nil {
		return nil
	}
	return o.output
}

// SetOutput sets the output writer.
//
// Parameters:
//   - `output`: the output writer
func (o *Options) SetOutput(output io.Writer) {
	if o == nil {
		return
	}
	o.output = output
}

// IsCallerReporter returns whether caller reporting is enabled.
//
// Returns:
//
// true if caller reporting is enabled.
func (o *Options) IsCallerReporter() bool {
	if o == nil {
		return false
	}
	return o.caller
}

// SetCallerReporter enables or disables caller reporting.
//
// Parameters:
//   - `enable`: whether to enable caller reporting
func (o *Options) SetCallerReporter(enable bool) {
	if o == nil {
		return
	}
	o.caller = enable
}

// GetCallerSkip returns the caller skip count.
//
// Returns:
//
// the number of stack frames to skip.
func (o *Options) GetCallerSkip() int {
	if o == nil {
		return 0
	}
	return o.callerSkip
}

// SetCallerSkip sets the caller skip count.
//
// Parameters:
//   - `skip`: the number of stack frames to skip
func (o *Options) SetCallerSkip(skip int) {
	if o == nil {
		return
	}
	o.callerSkip = skip
}

// GetFields returns a copy of the fields.
//
// Returns:
//
// a copy of the []Field slice configured in Options.
func (o *Options) GetFields() []Field {
	if o == nil || o.fields == nil {
		return nil
	}
	result := make([]Field, len(o.fields))
	copy(result, o.fields)
	return result
}

// SetFields sets the fields.
//
// Parameters:
//   - `fields`: the fields to set
func (o *Options) SetFields(fields []Field) {
	if o == nil {
		return
	}
	if fields == nil {
		o.fields = nil
		return
	}
	o.fields = make([]Field, len(fields))
	copy(o.fields, fields)
}

// AddFields appends fields to the existing field list.
//
// Parameters:
//   - `fields`: the fields to append
func (o *Options) AddFields(fields ...Field) {
	if o == nil {
		return
	}
	o.fields = append(o.fields, fields...)
}

// GetName returns the logger name.
//
// Returns:
//
// the name configured in Options.
func (o *Options) GetName() string {
	if o == nil {
		return ""
	}
	return o.name
}

// SetName sets the logger name.
//
// Parameters:
//   - `name`: the name to set
func (o *Options) SetName(name string) {
	if o == nil {
		return
	}
	o.name = name
}

// GetSamplingOpts returns the sampling options.
//
// Returns:
//
// the *SamplingOptions configured in Options.
func (o *Options) GetSamplingOpts() *SamplingOptions {
	if o == nil {
		return nil
	}
	return o.samplingOpts
}

// SetSamplingOpts sets the sampling options.
//
// Parameters:
//   - `opts`: the sampling options
func (o *Options) SetSamplingOpts(opts *SamplingOptions) {
	if o == nil {
		return
	}
	o.samplingOpts = opts
}

// GetRotationOpts returns the rotation options.
//
// Returns:
//
// the *RotationOptions configured in Options.
func (o *Options) GetRotationOpts() *RotationOptions {
	if o == nil {
		return nil
	}
	return o.rotationOpts
}

// SetRotationOpts sets the rotation options.
//
// Parameters:
//   - `opts`: the rotation options
func (o *Options) SetRotationOpts(opts *RotationOptions) {
	if o == nil {
		return
	}
	o.rotationOpts = opts
}

// ///////////////////////////////////////////////////////////////////////////
// SamplingOptions accessors
// ///////////////////////////////////////////////////////////////////////////

// GetFirst returns the number of messages always logged within the period.
//
// Returns:
//
// the first value.
func (s *SamplingOptions) GetFirst() int {
	if s == nil {
		return 0
	}
	return s.first
}

// SetFirst sets the number of messages always logged within the period.
//
// Parameters:
//   - `first`: the number of messages
func (s *SamplingOptions) SetFirst(first int) {
	if s == nil {
		return
	}
	s.first = first
}

// GetPeriod returns the sampling window duration.
//
// Returns:
//
// the period value.
func (s *SamplingOptions) GetPeriod() time.Duration {
	if s == nil {
		return 0
	}
	return s.period
}

// SetPeriod sets the sampling window duration.
//
// Parameters:
//   - `period`: the sampling window
func (s *SamplingOptions) SetPeriod(period time.Duration) {
	if s == nil {
		return
	}
	s.period = period
}

// GetThereafter returns the interval for logging after First is exhausted.
//
// Returns:
//
// the thereafter value.
func (s *SamplingOptions) GetThereafter() int {
	if s == nil {
		return 0
	}
	return s.thereafter
}

// SetThereafter sets the interval for logging after First is exhausted.
//
// Parameters:
//   - `thereafter`: every Nth message after First
func (s *SamplingOptions) SetThereafter(thereafter int) {
	if s == nil {
		return
	}
	s.thereafter = thereafter
}

// WithFirst sets the first count and returns the receiver for chaining.
//
// Parameters:
//   - `first`: the number of messages
//
// Returns:
//
// the receiver, for method chaining.
func (s *SamplingOptions) WithFirst(first int) *SamplingOptions {
	s.SetFirst(first)
	return s
}

// WithPeriod sets the sampling window duration and returns the receiver for chaining.
//
// Parameters:
//   - `period`: the sampling window
//
// Returns:
//
// the receiver, for method chaining.
func (s *SamplingOptions) WithPeriod(period time.Duration) *SamplingOptions {
	s.SetPeriod(period)
	return s
}

// WithThereafter sets the interval for logging after First is exhausted and returns the receiver for chaining.
//
// Parameters:
//   - `thereafter`: every Nth message after First
//
// Returns:
//
// the receiver, for method chaining.
func (s *SamplingOptions) WithThereafter(thereafter int) *SamplingOptions {
	s.SetThereafter(thereafter)
	return s
}

// ///////////////////////////////////////////////////////////////////////////
// SloggerConfig accessors
// ///////////////////////////////////////////////////////////////////////////

// GetLevel returns the log level string.
//
// Returns:
//
// the Level field value.
func (c *SloggerConfig) GetLevel() string {
	if c == nil {
		return ""
	}
	return c.Level
}

// SetLevel sets the log level string.
//
// Parameters:
//   - `level`: the log level string
func (c *SloggerConfig) SetLevel(level string) {
	if c == nil {
		return
	}
	c.Level = level
}

// GetFormatter returns the formatter type string.
//
// Returns:
//
// the Formatter field value.
func (c *SloggerConfig) GetFormatter() string {
	if c == nil {
		return ""
	}
	return c.Formatter
}

// SetFormatter sets the formatter type string.
//
// Parameters:
//   - `formatter`: the formatter type ("text" or "json")
func (c *SloggerConfig) SetFormatter(formatter string) {
	if c == nil {
		return
	}
	c.Formatter = formatter
}

// GetOutput returns the output configuration.
//
// Returns:
//
// the OutputConfig.
func (c *SloggerConfig) GetOutput() OutputConfig {
	if c == nil {
		return OutputConfig{}
	}
	return c.Output
}

// SetOutput sets the output configuration.
//
// Parameters:
//   - `output`: the output configuration
func (c *SloggerConfig) SetOutput(output OutputConfig) {
	if c == nil {
		return
	}
	c.Output = output
}

// GetFile returns the file configuration.
//
// Returns:
//
// the FileConfig.
func (c *SloggerConfig) GetFile() FileConfig {
	if c == nil {
		return FileConfig{}
	}
	return c.File
}

// SetFile sets the file configuration.
//
// Parameters:
//   - `file`: the file configuration
func (c *SloggerConfig) SetFile(file FileConfig) {
	if c == nil {
		return
	}
	c.File = file
}

// GetRotation returns the rotation configuration.
//
// Returns:
//
// the RotationConfig.
func (c *SloggerConfig) GetRotation() RotationConfig {
	if c == nil {
		return RotationConfig{}
	}
	return c.Rotation
}

// SetRotation sets the rotation configuration.
//
// Parameters:
//   - `rotation`: the rotation configuration
func (c *SloggerConfig) SetRotation(rotation RotationConfig) {
	if c == nil {
		return
	}
	c.Rotation = rotation
}

// GetArchive returns the archive configuration.
//
// Returns:
//
// the ArchiveConfig.
func (c *SloggerConfig) GetArchive() ArchiveConfig {
	if c == nil {
		return ArchiveConfig{}
	}
	return c.Archive
}

// SetArchive sets the archive configuration.
//
// Parameters:
//   - `archive`: the archive configuration
func (c *SloggerConfig) SetArchive(archive ArchiveConfig) {
	if c == nil {
		return
	}
	c.Archive = archive
}

// GetCaller returns the caller configuration.
//
// Returns:
//
// the CallerConfig.
func (c *SloggerConfig) GetCaller() CallerConfig {
	if c == nil {
		return CallerConfig{}
	}
	return c.Caller
}

// SetCaller sets the caller configuration.
//
// Parameters:
//   - `caller`: the caller configuration
func (c *SloggerConfig) SetCaller(caller CallerConfig) {
	if c == nil {
		return
	}
	c.Caller = caller
}

// GetColor returns the color configuration.
//
// Returns:
//
// the ColorConfig.
func (c *SloggerConfig) GetColor() ColorConfig {
	if c == nil {
		return ColorConfig{}
	}
	return c.Color
}

// SetColor sets the color configuration.
//
// Parameters:
//   - `color`: the color configuration
func (c *SloggerConfig) SetColor(color ColorConfig) {
	if c == nil {
		return
	}
	c.Color = color
}

// ///////////////////////////////////////////////////////////////////////////
// OutputConfig accessors
// ///////////////////////////////////////////////////////////////////////////

// IsConsole returns whether console output is enabled.
//
// Returns:
//
// the Console field value.
func (o *OutputConfig) IsConsole() bool {
	if o == nil {
		return false
	}
	return o.Console
}

// SetConsole enables or disables console output.
//
// Parameters:
//   - `console`: whether to enable console output
func (o *OutputConfig) SetConsole(console bool) {
	if o == nil {
		return
	}
	o.Console = console
}

// IsFile returns whether file output is enabled.
//
// Returns:
//
// the File field value.
func (o *OutputConfig) IsFile() bool {
	if o == nil {
		return false
	}
	return o.File
}

// SetFile enables or disables file output.
//
// Parameters:
//   - `file`: whether to enable file output
func (o *OutputConfig) SetFile(file bool) {
	if o == nil {
		return
	}
	o.File = file
}

// ///////////////////////////////////////////////////////////////////////////
// FileConfig accessors
// ///////////////////////////////////////////////////////////////////////////

// GetDirectory returns the log directory.
//
// Returns:
//
// the Directory field value.
func (f *FileConfig) GetDirectory() string {
	if f == nil {
		return ""
	}
	return f.Directory
}

// SetDirectory sets the log directory.
//
// Parameters:
//   - `directory`: the log directory path
func (f *FileConfig) SetDirectory(directory string) {
	if f == nil {
		return
	}
	f.Directory = directory
}

// GetInfoFile returns the info log file name.
//
// Returns:
//
// the InfoFile field value.
func (f *FileConfig) GetInfoFile() string {
	if f == nil {
		return ""
	}
	return f.InfoFile
}

// SetInfoFile sets the info log file name.
//
// Parameters:
//   - `infoFile`: the info log file name
func (f *FileConfig) SetInfoFile(infoFile string) {
	if f == nil {
		return
	}
	f.InfoFile = infoFile
}

// GetWarnFile returns the warn log file name.
//
// Returns:
//
// the WarnFile field value.
func (f *FileConfig) GetWarnFile() string {
	if f == nil {
		return ""
	}
	return f.WarnFile
}

// SetWarnFile sets the warn log file name.
//
// Parameters:
//   - `warnFile`: the warn log file name
func (f *FileConfig) SetWarnFile(warnFile string) {
	if f == nil {
		return
	}
	f.WarnFile = warnFile
}

// GetErrorFile returns the error log file name.
//
// Returns:
//
// the ErrorFile field value.
func (f *FileConfig) GetErrorFile() string {
	if f == nil {
		return ""
	}
	return f.ErrorFile
}

// SetErrorFile sets the error log file name.
//
// Parameters:
//   - `errorFile`: the error log file name
func (f *FileConfig) SetErrorFile(errorFile string) {
	if f == nil {
		return
	}
	f.ErrorFile = errorFile
}

// GetDebugFile returns the debug log file name.
//
// Returns:
//
// the DebugFile field value.
func (f *FileConfig) GetDebugFile() string {
	if f == nil {
		return ""
	}
	return f.DebugFile
}

// SetDebugFile sets the debug log file name.
//
// Parameters:
//   - `debugFile`: the debug log file name
func (f *FileConfig) SetDebugFile(debugFile string) {
	if f == nil {
		return
	}
	f.DebugFile = debugFile
}

// ///////////////////////////////////////////////////////////////////////////
// RotationConfig accessors
// ///////////////////////////////////////////////////////////////////////////

// Enabled returns whether rotation is enabled.
//
// Returns:
//
// the IsEnabled field value.
func (r *RotationConfig) Enabled() bool {
	if r == nil {
		return false
	}
	return r.IsEnabled
}

// SetEnabled enables or disables rotation.
//
// Parameters:
//   - `enabled`: whether to enable rotation
func (r *RotationConfig) SetEnabled(enabled bool) {
	if r == nil {
		return
	}
	r.IsEnabled = enabled
}

// GetMaxSizeMB returns the maximum file size in megabytes.
//
// Returns:
//
// the MaxSizeMB field value.
func (r *RotationConfig) GetMaxSizeMB() int64 {
	if r == nil {
		return 0
	}
	return r.MaxSizeMB
}

// SetMaxSizeMB sets the maximum file size in megabytes.
//
// Parameters:
//   - `maxSizeMB`: the maximum size in megabytes
func (r *RotationConfig) SetMaxSizeMB(maxSizeMB int64) {
	if r == nil {
		return
	}
	r.MaxSizeMB = maxSizeMB
}

// GetMaxAgeDays returns the maximum file age in days.
//
// Returns:
//
// the MaxAgeDays field value.
func (r *RotationConfig) GetMaxAgeDays() int {
	if r == nil {
		return 0
	}
	return r.MaxAgeDays
}

// SetMaxAgeDays sets the maximum file age in days.
//
// Parameters:
//   - `maxAgeDays`: the maximum age in days
func (r *RotationConfig) SetMaxAgeDays(maxAgeDays int) {
	if r == nil {
		return
	}
	r.MaxAgeDays = maxAgeDays
}

// IsCompress returns whether compression is enabled.
//
// Returns:
//
// the Compress field value.
func (r *RotationConfig) IsCompress() bool {
	if r == nil {
		return false
	}
	return r.Compress
}

// SetCompress enables or disables compression.
//
// Parameters:
//   - `compress`: whether to compress rotated files
func (r *RotationConfig) SetCompress(compress bool) {
	if r == nil {
		return
	}
	r.Compress = compress
}

// ///////////////////////////////////////////////////////////////////////////
// ArchiveConfig accessors
// ///////////////////////////////////////////////////////////////////////////

// Enabled returns whether archiving is enabled.
//
// Returns:
//
// the IsEnabled field value.
func (a *ArchiveConfig) Enabled() bool {
	if a == nil {
		return false
	}
	return a.IsEnabled
}

// SetEnabled enables or disables archiving.
//
// Parameters:
//   - `enabled`: whether to enable archiving
func (a *ArchiveConfig) SetEnabled(enabled bool) {
	if a == nil {
		return
	}
	a.IsEnabled = enabled
}

// GetPath returns the archive path.
//
// Returns:
//
// the Path field value.
func (a *ArchiveConfig) GetPath() string {
	if a == nil {
		return ""
	}
	return a.Path
}

// SetPath sets the archive path.
//
// Parameters:
//   - `path`: the archive path
func (a *ArchiveConfig) SetPath(path string) {
	if a == nil {
		return
	}
	a.Path = path
}

// GetFormat returns the date format for archive directories.
//
// Returns:
//
// the Format field value.
func (a *ArchiveConfig) GetFormat() string {
	if a == nil {
		return ""
	}
	return a.Format
}

// SetFormat sets the date format for archive directories.
//
// Parameters:
//   - `format`: the Go time layout string
func (a *ArchiveConfig) SetFormat(format string) {
	if a == nil {
		return
	}
	a.Format = format
}

// ///////////////////////////////////////////////////////////////////////////
// CallerConfig accessors
// ///////////////////////////////////////////////////////////////////////////

// Enabled returns whether caller reporting is enabled.
//
// Returns:
//
// the IsEnabled field value.
func (c *CallerConfig) Enabled() bool {
	if c == nil {
		return false
	}
	return c.IsEnabled
}

// SetEnabled enables or disables caller reporting.
//
// Parameters:
//   - `enabled`: whether to enable caller reporting
func (c *CallerConfig) SetEnabled(enabled bool) {
	if c == nil {
		return
	}
	c.IsEnabled = enabled
}

// ///////////////////////////////////////////////////////////////////////////
// ColorConfig accessors
// ///////////////////////////////////////////////////////////////////////////

// Enabled returns whether color output is enabled.
//
// Returns:
//
// the IsEnabled field value.
func (c *ColorConfig) Enabled() bool {
	if c == nil {
		return false
	}
	return c.IsEnabled
}

// SetEnabled enables or disables color output.
//
// Parameters:
//   - `enabled`: whether to enable color output
func (c *ColorConfig) SetEnabled(enabled bool) {
	if c == nil {
		return
	}
	c.IsEnabled = enabled
}
