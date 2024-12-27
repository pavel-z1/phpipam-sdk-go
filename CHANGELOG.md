# Changelog

Todas as modificações desse projeto devem ser adicionadas a esse arquivo.

O formato foi baseado em [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
e o projeto está seguindo o versionamento conforme [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [0.3.2] - 27/12/2024

### Fix

- Alterado IsGateway e ExcludePing para BoolInt

## [0.3.0] - 08/10/2024

### Fix

- Adicionado função GetSubnetsByVlan


## [0.2.3] - 17/07/2024

### Fix

- Ajustado Subnet.GatewayId de string para int

### Fix

- Removido anotação ",string" das structs com int
- Criado função BoolInt e substituido BoolIntString
- Ajustado testes unitários
