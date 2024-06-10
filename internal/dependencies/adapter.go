package dependencies

import (
	"github.com/warehouse/mail-service/internal/adapter/random"
	"github.com/warehouse/mail-service/internal/adapter/time"
)

func (d *dependencies) TimeAdapter() time.Adapter {
	if d.timeAdapter == nil {
		d.timeAdapter = time.NewAdapter(
			d.cfg.Time.Locale,
		)
	}

	return d.timeAdapter
}

func (d *dependencies) RandomAdapter() random.Adapter {
	if d.randomAdapter == nil {
		d.randomAdapter = random.NewAdapter()
	}
	return d.randomAdapter
}
