# Архитектура

```mermaid
graph BT
    %% Определение слоев
    subgraph Layer4 [Внешнее кольцо: ВХОДНЫЕ ТОЧКИ]
        direction LR
        Transport[transport/rest]
    end

    subgraph Layer3 [Инфраструктурное кольцо: РЕАЛИЗАЦИЯ]
        direction LR
        Postgres[repository/postgres]
        Redis[repository/redis]
    end

    subgraph Layer2 [Кольцо логики: ПРАВИЛА БИЗНЕСА]
        direction LR
        Service[service/logic]
    end

    subgraph Layer1 [ЯДРО: КОНТРАКТЫ И МОДЕЛИ]
        direction LR
        Domain{domain/core}
    end

    %% Направление зависимостей (все смотрят в центр/вниз)
    Transport --> Service
    Service --> Domain
    Postgres -.-> Domain
    Redis -.-> Domain
    
    %% Пояснение: Сервис использует интерфейсы репозитория из Domain
    Service --> |использует интерфейсы| Domain

    %% Стилизация
    style Layer1 fill:#e1f5fe,stroke:#01579b,stroke-width:2px
    style Layer2 fill:#fff9c4,stroke:#fbc02d,stroke-width:2px
    style Layer3 fill:#f1f8e9,stroke:#33691e,stroke-width:2px
    style Layer4 fill:#fce4ec,stroke:#880e4f,stroke-width:2px
    style Domain font-weight:bold,fill:#b3e5fc