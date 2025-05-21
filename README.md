# Pontos-Chave
- Garantir a disponibilidade de todos os postos na viagem
- Permitir planejamento e reserva de múltiplos postos a partir de qualquer servidor
- Utilizar requisições atômicas
- Comunicação entre servidores com MQTT e API REST
- "O cliente que reservar o primeiro ponto deve manter a prioridade na reserva sobre os trechos seguintes,
onde os demais clientes podem desistir ou continuar a compra da passagem escolhendo outros pontos de
carregamento disponíveis" ~Não entendi essa parte
- Pode usar framework
- Deve usar docker, API REST testada com Insominia ou Postman, MQTT e dados gerados aleatoriamente
- Entrega: 12/05

# Problema 
No problema anterior, foi desenvolvido um sistema inteligente de carregamento de veículos elétricos
que pode ser aplicado para gerenciar pontos de recarga em uma cidade. Neste problema, sua startup
identificou a dificuldade dos usuários do sistema em planejar e garantir as recargas necessárias para viagens
longas, entre cidades e estados. Em distâncias longas, é preciso ***garantir a disponibilidade sequencial*** para
completar a viagem dentro de um ***cronograma previsto***, com paradas planejadas de forma ***otimizada e segura***.

O novo desafio da sua equipe é aprimorar o sistema de recarga inteligente para 
***suportar o planejamento e a reserva antecipada de múltiplos pontos de recarga***, dentro de janelas de tempo definidas, ao
longo de uma rota específica entre cidades e estados. O objetivo é que, através de uma ***requisição atômica***, o
sistema possa ***consultar a disponibilidade e reservar*** uma sequência de pontos de recarga necessários para que
o veículo complete sua viagem sem o risco de ficar sem energia, evitando atrasos imprevistos devido à
indisponibilidade de carregadores. 

Para isso, é essencial que exista uma ***forma padronizada e coordenada de comunicação entre os servidores das empresas conveniadas envolvidas***. A ***comunicação entre os servidores*** deve ser realizada ***através de uma API*** projetada pela sua equipe de
desenvolvimento para permitir que um cliente possa, ***a partir de qualquer servidor***, ***reservar pontos*** de
carregamento disponíveis ***em diferentes empresas*** conveniadas seguindo as mesmas regras do sistema
centralizado original. 

Por exemplo, um cliente (carro) que está querendo viajar de João Pessoa à Feira de
Santana pode iniciar a requisição através do servidor da empresa A. Nesta requisição, o cliente escolhe um
ponto de carregamento entre João Pessoa e Maceió, da empresa A, outro ponto de carregamento entre
Maceió e Sergipe, da empresa B, e outro ponto de carregamento entre Sergipe a Feira de Santana, da empresa
C. O ***cliente que reservar o primeiro ponto*** deve manter a ***prioridade na reserva sobre os trechos seguintes***,
onde os ***demais clientes podem desistir ou continuar*** a compra da passagem escolhendo outros pontos de
carregamento disponíveis.

# Restrições
Diferente do anterior, neste problema é ***liberado o uso de frameworks*** de comunicação de terceiros para
implementar a solução do problema, limitados pelos seguintes requisitos:
- Para uma emulação realista do cenário proposto, os elementos da arquitetura devem ser executados
em ***contêineres Docker***, ***executados em computadores distintos*** no laboratório;
- A interface entre os servidores deve ser projetada e implementada através de 
***protocolo baseado em API REST***, podendo ser ***testada*** na apresentação ***através de*** softwares como ***Insomnia ou Postman***;
- Os ***carros*** devem ser ***simulados*** através de um software para geração de dados fictícios, onde os ***dados***
devem ser ***gerados aleatoriamente*** passando a tendência da ***descarga da bateria (rápida, lenta, etc.)***;
- Na comunicação dos carros com o servidor, ***ao invés de uma API de sockets***, estabeleceu-se que a
solução deve adotar o ***padrão usado na Internet das Coisas (IoT)***, com o ***protocolo Message Queue Telemetry Transport (MQTT)***,
classificado como um protocolo ***Machine-to-Machine (M2M)***.

# Cronograma

Entrega: 12/05
Entrega fora do prazo: -20% da nota e -5% por dia de atraso

Apresentação: 12/05 e 14/05

# Avaliação

A nota final será composta por três critérios de avaliação:
1. Desempenho individual (25%)
2. Documentação (25%)
3. Produto Final (código incluso) (50%)



# Arquitetura da Solução

O sistema foi desenvolvido com uma **arquitetura distribuída baseada em microserviços**, composta pelos seguintes componentes principais:

## 🧩 Componentes Principais

### 📡 Servidor (API REST + MQTT)
Responsável por:
- Gerenciar as estações de recarga, rotas e reservas.
- Cada servidor representa uma empresa distinta.
- Expõe endpoints **REST** para comunicação entre servidores.
- Integra-se ao **broker MQTT** para comunicação com os clientes (carros e estações).

### 🚗 Clientes (Carros e Estações)
Simulam:
- **Usuários (carros)** e **pontos de recarga (estações)**.
- Comunicam-se com o servidor via **MQTT**, publicando e recebendo mensagens em tópicos específicos.
- Operações suportadas:
  - Consulta de rotas
  - Reserva de estações
  - Liberação de estações

### 🔀 Broker MQTT
- Atua como **middleware para troca de mensagens assíncronas** entre clientes e servidores.
- Permite **desacoplamento** entre componentes e promove **escalabilidade**.

### 🗃️ Banco de Dados
- Responsável pela **persistência de informações**:
  - Estações
  - Rotas
  - Reservas
  - Identificação dos servidores

---

## 🏗️ Classificação da Arquitetura

A solução é classificada como uma:

### 👉 Arquitetura de Microserviços Distribuídos
- **Orientada a eventos** (via MQTT)
- **Requisições síncronas** (via REST)

Cada componente possui:
- **Responsabilidades bem definidas**
- **Comunicação padronizada**

### ✅ Benefícios:
- Escalabilidade
- Modularidade
- Facilidade de manutenção


## 📡 Protocolo de Comunicação

A solução utiliza **dois protocolos principais** para a comunicação entre os componentes do sistema:

---

### 1. 🛰️ MQTT (Message Queue Telemetry Transport)

**Comunicação entre:**  
➡️ Clientes (**carros/estações**) e **servidores**  
**Tipo:** Comunicação **assíncrona**, **orientada a eventos**

#### 📌 Tópicos e Payloads

**Exemplos de tópicos:**
- `car/{serverIP}/{carID}/consult` — Consulta de rotas
- `car/{serverIP}/{carID}/reserve` — Reserva de rota
- `car/{serverIP}/{carID}/finishroute` — Finalização de rota
- `station/{serverIP}/{stationID}/birth` — Nascimento de estação

**Formato dos payloads:**  
- JSON (ex: `SelectRouteMessage`, `FinishRouteMessage`, `CarInfo`, etc.)

#### 🔁 Fluxo típico:
1. O carro publica uma **mensagem de reserva** em um tópico MQTT.
2. O servidor processa a requisição.
3. O servidor responde em um **tópico específico de resposta** para aquele carro.
4. Mensagens de **finalização de rota** seguem o mesmo padrão.

---

### 2. 🌐 API REST

**Comunicação entre:**  
➡️ **Servidores**  
**Tipo:** Comunicação **síncrona**, baseada em **requisições HTTP**

#### 📌 Principais endpoints:

- `PUT /server/:sid/stations/:id/prepare`  
  ➤ Prepara uma estação remota para reserva

- `PUT /server/:sid/stations/:id/commit`  
  ➤ Efetiva a reserva remota

- `PUT /server/:sid/stations/:id/release`  
  ➤ Libera uma estação remota

- `POST /servers/register`  
  ➤ Registra um novo servidor

#### 📎 Parâmetros:
- IDs de estação, IDs de servidor, IDs de carro
- Payloads em **JSON** (ex: `{ "car_id": 123 }`)

#### 📤 Retornos:
- Status HTTP (`200 OK`, `400 Bad Request`, etc.)
- Mensagens JSON indicando **sucesso ou erro**

---

### 🧭 Resumo do Fluxo de Reserva de Rota

1. 🚗 O carro publica uma **mensagem de reserva via MQTT**.
2. 🧠 O servidor consulta e reserva **estações locais e remotas** via **API REST**.
3. 🔁 O servidor executa o protocolo **2PC (Prepare/Commit)** entre servidores via REST.
4. 📩 O servidor responde ao carro via MQTT com o resultado da operação.

## 🚦 Roteamento

O sistema implementa um **roteamento distribuído** para calcular e apresentar ao usuário todas as rotas possíveis entre origem e destino, considerando os pontos de recarga disponíveis em servidores de todas as companhias.

### 🧮 Como funciona o cálculo de rotas?

- **Consulta de rotas:**  
  O cliente (carro) envia uma mensagem de consulta via MQTT para o servidor, informando as cidades de origem e destino.
- **Busca distribuída:**  
  O servidor consulta seu banco de dados por todas as rotas possíveis entre as cidades informadas, levando em conta as estações disponíveis (ativas) em sua própria empresa e, se necessário, consulta outros servidores para incluir estações de outras companhias.
- **Resposta ao usuário:**  
  O servidor retorna ao cliente todas as rotas possíveis, cada uma composta por uma sequência de estações de recarga (de diferentes empresas, se necessário), garantindo que o usuário possa planejar a viagem completa sem risco de ficar sem energia.

### 📋 Exemplo de fluxo

1. 🚗 O carro consulta rotas de João Pessoa para Feira de Santana.
2. 🧠 O servidor retorna múltiplas opções de rotas, cada uma com diferentes pontos de recarga (ex: João Pessoa → Maceió [Empresa A], Maceió → Sergipe [Empresa B], Sergipe → Feira de Santana [Empresa C]).
3. 👤 O usuário escolhe a rota desejada e inicia o processo de reserva.

### ✅ O sistema garante:

- **Cálculo distribuído:** As rotas podem envolver estações de várias empresas, consultando diferentes servidores.
- **Exibição de todas as possibilidades:** O usuário visualiza todas as rotas possíveis, considerando a disponibilidade dos pontos de recarga em todos os servidores participantes.
- **Reserva atômica:** A reserva dos pontos de recarga ao longo da rota é feita de forma coordenada, garantindo que o usuário só inicie a viagem se todos os pontos estiverem disponíveis.

## 🤝 Concorrência Distribuída

Para evitar que o mesmo ponto de recarga seja reservado por clientes distintos no mesmo horário, o sistema emprega o protocolo de commit em duas fases (2PC):

- **Fase de preparação:** Cada estação envolvida na rota recebe uma requisição de "prepare" e só aceita se estiver realmente disponível.
- **Fase de commit:** Se todas as estações confirmarem a preparação, a reserva é efetivada em todas. Se alguma não puder reservar, todas as reservas são abortadas.
- **Garantia:** Nenhum ponto é reservado simultaneamente para clientes diferentes, mesmo em ambiente distribuído e com múltiplos servidores/empresas.

Esse controle é feito de forma distribuída, coordenando as reservas entre servidores via API REST e garantindo a consistência do sistema.

## 🔒 Confiabilidade da Solução

O sistema foi projetado para garantir a confiabilidade e a consistência das reservas mesmo diante de falhas de comunicação ou desconexão temporária de servidores das companhias.

- **Protocolo 2PC (Two-Phase Commit):**  
  Utiliza o protocolo de commit em duas fases para garantir que uma reserva só será efetivada se todos os servidores participantes confirmarem a operação. Caso algum servidor fique indisponível durante o processo, a reserva é automaticamente abortada para todos, evitando inconsistências.

- **Persistência:**  
  O estado das reservas e estações é salvo no banco de dados, permitindo que servidores retomem o processamento corretamente após uma falha ou reconexão.

- **Recuperação de falhas:**  
  Se um servidor desconectar durante uma reserva, o sistema aborta a operação e libera os recursos envolvidos. Ao reconectar, o servidor pode consultar o banco de dados para retomar seu estado.

- **Garantia de concorrência distribuída:**  
  Mesmo em cenários de falha, o sistema impede que dois clientes reservem o mesmo ponto de recarga no mesmo horário, mantendo a integridade e a atomicidade das operações distribuídas.

Dessa forma, o sistema continua garantindo a concorrência distribuída e a finalização correta das reservas, mesmo com desconexão e reconexão dos servidores das companhias. 
