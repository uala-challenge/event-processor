# Event Processor

## Descripción
El servicio **Event Processor** es una aplicación que procesa mensajes del tópico en batch, asegurando la correcta gestión y distribución de los eventos generados en la plataforma.

## Características principales
- Consumo eficiente de eventos desde el tópico.
- Procesamiento batch para optimización del rendimiento.
- Uso de concurrencia con patrones **fan-in/fan-out** y **worker pools** para mejorar la eficiencia del procesamiento.
- Manejo de errores y reintentos automáticos.
- Arquitectura escalable y desacoplada.
- Logs detallados para monitoreo y depuración.

## Tecnologías utilizadas
- **Golang**: Lenguaje de desarrollo principal.
- **SQS**: Broker de mensajería (dependiendo de la configuración).
- **Docker**: Contenedorización para facilitar despliegue.
- **Logrus**: Manejo de logs estructurados.

## Instalación y configuración
### Prerrequisitos
- Tener instalado **Go** (versión 1.24).
- Contar con **Docker** y **Docker Compose** (para pruebas locales).
- Configurar las variables de entorno adecuadas.

### Clonar el repositorio
```bash
  git clone https://github.com/uala-challenge/event-processor.git
  cd event-processor
```

### Configuración
El servicio utiliza un archivo de configuración en formato YAML. Antes de ejecutar, asegúrate de configurar correctamente `config.yaml`:
```yaml
broker:
  type: kafka  # Opciones: kafka, rabbitmq
  url: "localhost:9092"
  topic: "events"
  group_id: "event-processor-group"

batch:
  size: 100
  timeout: 5s

logging:
  level: "info"
```

### Ejecutar en entorno local
#### Usando Go directamente
```bash
  go run main.go
```

#### Usando Docker
```bash
  docker-compose up --build
```

## Uso
El Event Processor consume eventos del tópico especificado y los procesa en lotes de acuerdo con la configuración. Se hace uso de **fan-in/fan-out** y **worker pools** para procesar los eventos concurrentemente, asegurando una distribución eficiente de la carga y optimizando el rendimiento.


## Testing
Para ejecutar las pruebas unitarias:
```bash
  go test ./...
```
