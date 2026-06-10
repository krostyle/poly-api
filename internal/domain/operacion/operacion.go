package operacion

import "time"

// Operacion representa una transacción impugnada dentro de un caso.
type Operacion struct {
	ID        string
	CasoID    string
	MedioPago string
	Relacion  string
	MontoCLP  int64
	MontoUF   *float64
	FechaOp   time.Time
	CreatedAt time.Time
}
