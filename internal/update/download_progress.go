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
	if total <= 0 {
		return base
	}
	pct := base + int(done*int64(span)/total)
	if pct > base+span {
		return base + span
	}
	return pct
}

func formatDownloadMessage(done, total int64) string {
	if total <= 0 {
		return "Скачивание обновления…"
	}
	pct := done * 100 / total
	if pct > 100 {
		pct = 100
	}
	return fmt.Sprintf("Скачивание… %d%%", pct)
}

func copyWithProgress(dst io.Writer, src io.Reader, total int64, onChange func(downloaded, total int64)) (int64, error) {
	if onChange != nil && total > 0 {
		onChange(0, total)
	}
	if total > 0 && onChange != nil {
		pw := &progressWriter{total: total, onChange: onChange}
		return io.Copy(io.MultiWriter(dst, pw), src)
	}
	return io.Copy(dst, src)
}
