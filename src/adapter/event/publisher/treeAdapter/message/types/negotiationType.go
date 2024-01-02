package types

type NegotiationType = string

const (
	NegotiationTypeMesa          NegotiationType = "MESA"
	NegotiationTypeCorretora     NegotiationType = "CORRETORA"
	NegotiationTypeMesaCorretora NegotiationType = "MESA_CORRETORA"
)