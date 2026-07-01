package update

import (
	"fmt"
	"io"
)

type progressWriter struct {
	total    int64
	written  int64
	onChange func(downloaded, total int64)
}

func (p *progressWriter) Write(b []byte) (int, error) {
	n := len(b)
	p.written += int64(n)
	if p.onChange != nil {
		p.onChange(p.written, p.total)
	}
	return n, nil
}

func downloadPercent(done, total int64, base, span int) int {
	if total > 0 {
		pct := base + int(done*int64(span)/total)
		if pct > base+span {
			return base + span
		}
		return pct
	}
	if done <= 0 {
		return base
	}
	// Content-Length unknown (redirect/chunked): creep toward span so the bar still moves.
	cap := base + span - 1
	if cap <= base {
		return base
	}
	step := int64(span - 1)
	if done < 512*1024 {
		return base + int(done*step/(512*1024))
	}
	return cap
}

func formatDownloadMessage(done, total int64) string {
	if total <= 0 {
		if done <= 0 {
			return "Скачивание обновления…"
		}
		return fmt.Sprintf("Скачивание… %s", formatByteCount(done))
	}
	pct := done * 100 / total
	if pct > 100 {
		pct = 100
	}
	return fmt.Sprintf("Скачивание… %d%%", pct)
}

func formatByteCount(n int64) string {
	const unit = 1024
	if n < unit {
		return fmt.Sprintf("%d B", n)
	}
	div, exp := int64(unit), 0
	for size := n / unit; size >= unit; size /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %cB", float64(n)/float64(div), "KMGTPE"[exp])
}

func copyWithProgress(dst io.Writer, src io.Reader, total int64, onChange func(downloaded, total int64)) (int64, error) {
	if onChange != nil {
		onChange(0, total)
		pw := &progressWriter{total: total, onChange: onChange}
		return io.Copy(io.MultiWriter(dst, pw), src)
	}
	return io.Copy(dst, src)
}

// FormatDownloadMessage formats a user-facing download status line.
func FormatDownloadMessage(done, total int64) string {
	return formatDownloadMessage(done, total)
}

// DownloadPercent maps byte progress into a UI percent range.
func DownloadPercent(done, total int64, base, span int) int {
	return downloadPercent(done, total, base, span)
}
