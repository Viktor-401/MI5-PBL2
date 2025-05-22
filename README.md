# MI5-PBL2

Este projeto é uma aplicação desenvolvida para o Problema 2 da Disciplina TEC502 - Concorrência e Conectividade.

## Como rodar

### Pré-requisitos

- Docker instalado

## Ordem de execução dos componentes

Para garantir o funcionamento correto do sistema, inicie os componentes na seguinte ordem:

1. **Broker MQTT**
2. **MongoDB**
3. **API (Server)**
4. **Clientes** (car, station)

### Configuração necessária para o Broker MQTT

Antes de rodar o broker, é necessário configurar os seguintes arquivos:

- `clients/broker/mosquitto.conf`: define a porta e permissões do broker.
- `clients/types/config.go`: define o IP e porta do broker para os clientes.
- `api/types/config.go`: define o IP e porta do broker para a API.

Certifique-se de que o IP e a porta definidos nesses arquivos estejam corretos e compatíveis com o ambiente onde o broker será executado.

### Rodando em modo **host**

No modo host, utilize os comandos abaixo para cada componente, seguindo a ordem acima:

#### 1. Broker MQTT

```bash
cd clients
make hostbroker
```

#### 2. MongoDB

```bash
cd api
make hostmongo
```

#### 3. API

```bash
cd api
make hostserver NAME=api PORT=8080 DBIP=127.0.0.1
```
> **Nota:** Defina `DBIP` como o IP do host onde o MongoDB está rodando.

#### 4. Clients

```bash
cd clients
make hostcar PORT=8081
make hoststation PORT=8082
```

### Rodando em modo **bridge** (rede bridge/app-net)

No modo bridge, utilize os comandos abaixo para cada componente, seguindo a ordem acima:

#### 1. Broker MQTT

```bash
cd clients
make broker
```

#### 2. MongoDB

```bash
cd api
make mongo
```

#### 3. API

```bash
cd api
make server NAME=api PORT=8080
```

#### 4. Clients

```bash
cd clients
make car PORT=8081
make station PORT=8082
```

### Sobre as variáveis NAME e PORT

- `NAME`: Define o nome do container da API/Server (nome como é armazenado no db).
- `PORT`: Define a porta exposta pelo serviço (por exemplo, 8080 para API, 8081 para car, 8082 para station).

### Observações


