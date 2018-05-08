package soap

import (
	"github.com/fiorix/wsdl2go/soap"
)

// Namespace was auto-generated from WSDL.
var Namespace = "http://131.1.18.108:8084/namespace/default"

// NewServiciosWebRPC creates an initializes a ServiciosWebRPC.
func NewServiciosWebRPC(cli *soap.Client) ServiciosWebRPC {
	return &serviciosWebRPC{cli}
}

// ServiciosWebRPC was auto-generated from WSDL
// and defines interface for the remote service. Useful for testing.
type ServiciosWebRPC interface {
	// WsAlertaMuestrasInactivas was auto-generated from WSDL.
	WsAlertaMuestrasInactivas(wsNumRecepcion string, wsCodEquipo string, wsCodInactivos string) (string, error)

	// WsCrearReto was auto-generated from WSDL.
	WsCrearReto(wDescripcionReto string, wResponsableReto string, wFechaReto Date, wHoraReto Time, wCorreosAdicionalesReto string, wOTReto string, wRetoXModificacion bool) error
}

// Date in WSDL format.
type Date string

// Time in WSDL format.
type Time string

// Operation wrapper for WsAlertaMuestrasInactivas.
// OperationWsAlertaMuestrasInactivasRequest was auto-generated
// from WSDL.
type OperationWsAlertaMuestrasInactivasRequest struct {
	WsNumRecepcion *string `xml:"wsNumRecepcion,omitempty" json:"wsNumRecepcion,omitempty" yaml:"wsNumRecepcion,omitempty"`
	WsCodEquipo    *string `xml:"wsCodEquipo,omitempty" json:"wsCodEquipo,omitempty" yaml:"wsCodEquipo,omitempty"`
	WsCodInactivos *string `xml:"wsCodInactivos,omitempty" json:"wsCodInactivos,omitempty" yaml:"wsCodInactivos,omitempty"`
}

// Operation wrapper for WsAlertaMuestrasInactivas.
// OperationWsAlertaMuestrasInactivasResponse was auto-generated
// from WSDL.
type OperationWsAlertaMuestrasInactivasResponse struct {
	WsResultado *string `xml:"wsResultado,omitempty" json:"wsResultado,omitempty" yaml:"wsResultado,omitempty"`
}

// Operation wrapper for WsCrearReto.
// OperationWsCrearRetoRequest was auto-generated from WSDL.
type OperationWsCrearRetoRequest struct {
	WDescripcionReto        *string `xml:"wDescripcionReto,omitempty" json:"wDescripcionReto,omitempty" yaml:"wDescripcionReto,omitempty"`
	WResponsableReto        *string `xml:"wResponsableReto,omitempty" json:"wResponsableReto,omitempty" yaml:"wResponsableReto,omitempty"`
	WFechaReto              *Date   `xml:"wFechaReto,omitempty" json:"wFechaReto,omitempty" yaml:"wFechaReto,omitempty"`
	WHoraReto               *Time   `xml:"wHoraReto,omitempty" json:"wHoraReto,omitempty" yaml:"wHoraReto,omitempty"`
	WCorreosAdicionalesReto *string `xml:"wCorreosAdicionalesReto,omitempty" json:"wCorreosAdicionalesReto,omitempty" yaml:"wCorreosAdicionalesReto,omitempty"`
	WOTReto                 *string `xml:"wOTReto,omitempty" json:"wOTReto,omitempty" yaml:"wOTReto,omitempty"`
	WRetoXModificacion      *bool   `xml:"wRetoXModificacion,omitempty" json:"wRetoXModificacion,omitempty" yaml:"wRetoXModificacion,omitempty"`
}

// Operation wrapper for WsCrearReto.
// OperationWsCrearRetoResponse was auto-generated from WSDL.
type OperationWsCrearRetoResponse struct {
}

// serviciosWebRPC implements the ServiciosWebRPC interface.
type serviciosWebRPC struct {
	cli *soap.Client
}

// WsAlertaMuestrasInactivas was auto-generated from WSDL.
func (p *serviciosWebRPC) WsAlertaMuestrasInactivas(wsNumRecepcion string, wsCodEquipo string, wsCodInactivos string) (string, error) {
	α := struct {
		M OperationWsAlertaMuestrasInactivasRequest `xml:"tns:wsAlertaMuestrasInactivas"`
	}{
		OperationWsAlertaMuestrasInactivasRequest{
			&wsNumRecepcion,
			&wsCodEquipo,
			&wsCodInactivos,
		},
	}

	γ := struct {
		M OperationWsAlertaMuestrasInactivasResponse `xml:"wsAlertaMuestrasInactivasResponse"`
	}{}
	if err := p.cli.RoundTripWithAction("ServiciosWeb#wsAlertaMuestrasInactivas", α, &γ); err != nil {
		return "", err
	}
	return *γ.M.WsResultado, nil
}

// WsCrearReto was auto-generated from WSDL.
func (p *serviciosWebRPC) WsCrearReto(wDescripcionReto string, wResponsableReto string, wFechaReto Date, wHoraReto Time, wCorreosAdicionalesReto string, wOTReto string, wRetoXModificacion bool) error {
	α := struct {
		M OperationWsCrearRetoRequest `xml:"tns:wsCrearReto"`
	}{
		OperationWsCrearRetoRequest{
			&wDescripcionReto,
			&wResponsableReto,
			&wFechaReto,
			&wHoraReto,
			&wCorreosAdicionalesReto,
			&wOTReto,
			&wRetoXModificacion,
		},
	}

	γ := struct {
		M OperationWsCrearRetoResponse `xml:"wsCrearRetoResponse"`
	}{}
	if err := p.cli.RoundTripWithAction("ServiciosWeb#wsCrearReto", α, &γ); err != nil {
		return err
	}
	return nil
}
