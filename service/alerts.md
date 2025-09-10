# Alertas y Warinings

----

## A001-A099/W001-W099-> Alertas/Warnings de LTM 
## A101-A199/W101-W199-> Alertas/Warnings de LTA
## A201-A299/W201-W299-> Alertas/Warnings de LTC 
## A301-A399/W301-W399-> Alertas/Warnings de SCR 

----

# LTM

## ⚠️Warning LTM: W001 --> Variación de cobertura

**Variación de cobertura de +- 25%, y nunca descendiendo por debajo del 50%.**

- **Descripción**: Posible instalación del LTM240 en techo de cabina
- **Recomendación**: ??

----

## ‼️Alerta LTM: A001

**"signal": "30%"**

- **Descripción**: Cobertura del 30% o menor. Cobertura baja.
- **Recomendación**: ??

----

## ‼️Alerta LTM: A002

**Variación de cobertura de +- 25%, y con algún valor por debajo del 50%.**

- **Descripción**: Incidencia de cobertura. LTM240 instalado en localización desfavorable (posible techo de cabina).
- **Recomendación**:  ??

----

## ‼️Alerta LTM: A003

**Cambio habitual en disp "Ctfs": si existen variaciones constantes en los eventos recibidos es una alerta.**

- **Descripción**: Ejemplo: Evento a continuación: Ctfs. Vemos una variacion en la detección de cabina. Módulo mal conectado o en mal estado. 
- **Recomendación**:  Ver histórico.

----

## ‼️Alerta LTM: A004

**"vusb":"4.97":El voltaje usb normal debe estar entre 4.6 y 5.60**

- **Descripción**:. Problema de alimentación. 
- **Recomendación**: ??

----

## ‼️Alerta LTM: A005 - Mode Flipping

**"vbat":"4.12": Esta alerta tiene que ver con los eventos recibidos.** (Bucle MODO_RED -> MODO_BATERIA)

- **Descripción**: Si existe una acumulación superior a 10 de cada uno de los eventos EV:4 y EV:5, es: Problema de alimentación y/o batería. (Cache alerta por IP).
- **Recomendación**:  ??
