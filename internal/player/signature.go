package player

import (
	"context"
	"fmt"

	playwright "github.com/mxschmitt/playwright-go"
)

// Normalized signature stroke inside the element (after a 10% inset).
var signatureStroke = [][2]float64{
	{0.05, 0.55},
	{0.45, 0.82},
	{0.95, 0.38},
}

const signatureInset = 0.1

func drawSignature(ctx context.Context, page playwright.Page, selector string) error {
	locator := page.Locator(selector)
	if err := waitForLocator(ctx, locator, playwright.LocatorWaitForOptions{
		State: playwright.WaitForSelectorStateVisible,
	}); err != nil {
		return fmt.Errorf("draw signature failed: %w", err)
	}
	if err := locator.ScrollIntoViewIfNeeded(); err != nil {
		return fmt.Errorf("draw signature scroll: %w", err)
	}
	if _, err := locator.Evaluate(signatureDispatchJS, nil); err == nil {
		return nil
	}
	return drawSignatureWithMouse(locator)
}

// signatureDispatchJS draws using getBoundingClientRect() at action time (responsive-safe).
const signatureDispatchJS = `(el) => {
  const points = [[0.05, 0.55], [0.45, 0.82], [0.95, 0.38]];
  const rect = el.getBoundingClientRect();
  if (!rect || rect.width < 4 || rect.height < 4) {
    throw new Error('signature target too small');
  }
  const inset = 0.1;
  const span = 1 - 2 * inset;
  const toXY = (nx, ny) => ({
    x: rect.left + rect.width * (inset + nx * span),
    y: rect.top + rect.height * (inset + ny * span),
  });
  const emit = (type, x, y) => {
    const base = { bubbles: true, cancelable: true, clientX: x, clientY: y, view: window };
    el.dispatchEvent(new MouseEvent(type, base));
    if (typeof PointerEvent !== 'undefined') {
      el.dispatchEvent(new PointerEvent(type.replace('mouse', 'pointer'), {
        ...base,
        pointerId: 1,
        pointerType: 'mouse',
        isPrimary: true,
      }));
    }
  };
  const first = toXY(points[0][0], points[0][1]);
  emit('mousedown', first.x, first.y);
  for (let i = 1; i < points.length; i++) {
    const p = toXY(points[i][0], points[i][1]);
    emit('mousemove', p.x, p.y);
  }
  const last = toXY(points[points.length - 1][0], points[points.length - 1][1]);
  emit('mouseup', last.x, last.y);
}`

func drawSignatureWithMouse(locator playwright.Locator) error {
	box, err := locator.BoundingBox()
	if err != nil || box == nil {
		return fmt.Errorf("draw signature: element has no bounding box")
	}
	startX, startY := signaturePoint(box, signatureStroke[0][0], signatureStroke[0][1])
	midX, midY := signaturePoint(box, signatureStroke[1][0], signatureStroke[1][1])
	endX, endY := signaturePoint(box, signatureStroke[2][0], signatureStroke[2][1])
	page, err := locator.Page()
	if err != nil {
		return fmt.Errorf("draw signature page: %w", err)
	}
	mouse := page.Mouse()
	if err := mouse.Move(startX, startY); err != nil {
		return err
	}
	if err := mouse.Down(); err != nil {
		return err
	}
	if err := mouse.Move(midX, midY); err != nil {
		return err
	}
	if err := mouse.Move(endX, endY); err != nil {
		return err
	}
	return mouse.Up()
}

func signaturePoint(box *playwright.Rect, nx, ny float64) (float64, float64) {
	span := 1 - 2*signatureInset
	return box.X + box.Width*(signatureInset+nx*span), box.Y + box.Height*(signatureInset+ny*span)
}
