package hass

type Hass struct {
	bearer   string
	endpoint string
	entity   string
}

func NewHass(
	bearer string,
	endpoint string,
	entity string,
) *Hass {
	return &Hass{
		bearer:   bearer,
		endpoint: endpoint,
		entity:   entity,
	}
}
