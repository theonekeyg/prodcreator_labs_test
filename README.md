# Description
Монорепа состоит из 3 компонентов:
* `aggspotter_keys` - Ключи для демо
* `solana_solution` - Сабрепа для контрактов
* `go_aggspotter` - Golang CLI с байднингами для контракта для запуска демо-скрипта

# Demo
Что-бы запустить демо, нужно иметь иметь локально установленную
[солану](https://docs.solana.com/cli/install-solana-cli-tools),
[anchor](https://www.anchor-lang.com/docs/installation),
[golang](https://go.dev/).

## Step 1. Modify and Build contracts

Для начала нужно изменить `OWNER_GUARD` в
`./solana_solution/programs/AggregationSpotter/src/lib.rs` на ключ, который
будет использоваться при деплоя (системный солана ключ, обычно
`~/.config/solana/id.json`), это защитник от ботов, защищает `initialize` метод
сразу после деплоя.

Далее сбилдим контракт.

`$ cd solana_solution && anchor build`

Теперь нужно будет изменить авто-сгенерированные ключи контрактов на
подготовленные мною, так как скрипт в golang будет ожидать несколько констант,
плюс не нужно будет менять id в контрактах.

`$ cp ../aggspotter_keys/aggregation_spotter-keypair.json ./target/deploy/`

`$ cp ../aggspotter_keys/test_linker-keypair.json ./target/deploy/`

Далее ребилдим контракт уже с правильными ключами и запускаем тесты.

`$ anchor build && anchor test`

## Step 2. Run local node and deploy program

Откроем новый терминал и зайдем в временную папку, и запустим в ней локальную
ноду.

`$ solana-test-validator`

Для тестов эта нода поднималась у нас автоматически, изменим этот параметр
cluster в `Anchor.toml`, поставим `http://127.0.0.1:8899` что-бы коннектится к
запущенной ноде.

В старом терминале задеплоим контракт и запустим migrate скрипт

`$ anchor deploy && anchor migrate`

## Step 3. Run golang test-script

Запустим golang скрипт, который создает новую операцию on-chain, голосует за нее
с указанных киперов, и при подтверждении запускает ее через контракт.

Зайдем в `go_aggspotter`

`$ cd ../go_aggspotter`

Загрузим зависимости

`$ go get`

Запустим подготовленный скрипт с `LOG_LEVEL` переменной среды на DEBUG

`$ LOG_LEVEL=DEBUG ./run.sh`

В логах будет выведена информация о транзакции, логах, аккаунтах, etc.

# Задача 

Задача:
Для блокчейна Solana на языке Rust переписать контракт AggregationSpotter.
Написать контракт TestLinker, с функцией testLink(args) которая будет
взаимодействовать с любым, на выбор кандидата, контрактом из тестнета Solana и
делать изменение его стейта (это может быть стейкинг контракт, dex или любой
другой). Суть в том, что контракт может находиться в другом шарде и вызов к нему
должен быть ассинхронным. После асинхронного выполнения операции на стороннем
контракте - выбросить ивент Ok(“ok).

Написать на языке Go программу-скрипт, которая с 3х разных аккаунтов - киперов
предложит и подтвердит операцию в контракте AggregationSpotter
(proposeOperation) на контракт TestLinker и функцию testLink. После того, как
операция была подтверждена (все киперы отправили предложение), отловить ивент
ProposalApproved(opHash) и выполнить операцию executeOperation. После, программа
должна отловить ивент Ok и вывести его.
