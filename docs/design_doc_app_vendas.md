# Design Document - Sistema de Controle de Vendas para Profissionais Autônomos

## 1. Objetivo do Sistema

O objetivo principal do sistema é tornar fácil o controle das vendas para profissionais autônomos e pequenas empresas que ainda não precisam de uma solução mais robusta. O sistema oferecerá uma experiência simples e intuitiva para que esses profissionais possam registrar, acompanhar e organizar suas vendas e recebimentos de forma eficiente.

## 2. Usuários Principais

O sistema será utilizado exclusivamente por **profissionais autônomos** e **pequenos empreendedores**.

- Haverá **apenas um tipo de usuário**, que será o próprio profissional.
- Este usuário terá acesso a todas as funcionalidades.
- O sistema é focado em uso individual, sem suporte a múltiplos perfis ou permissões diferenciadas.

## 3. Funcionalidades Principais

- **Cadastro de Clientes**  
  - Nome, telefone e observações adicionais.
  - Cadastro de múltiplos endereços e contatos por cliente.

- **Registro de Vendas**  
  - Registro de itens vendidos, quantidade, valor unitário e total.
  - Data da venda e forma de pagamento.
  - Status de pagamento: "A pagar" ou "Pago".

- **Parcelamento de Vendas**  
  - Opção para dividir vendas em parcelas.
  - Cálculo automático dos vencimentos.

- **Lembretes de Cobrança**  
  - Notificações visuais no app sobre cobranças vencidas ou próximas do vencimento.

- **Relatórios de Vendas**  
  - Filtros por data, cliente ou status de pagamento.
  - Apresentação simples com foco em clareza e acessibilidade.

## 4. Plataforma e Stack Tecnológica

- **Plataforma**: Aplicativo mobile (Android e iOS).
- **Frontend (App)**: Desenvolvido em **Flutter** (multiplataforma). A escolha se baseia na experiência prévia do desenvolvedor e na performance do framework, que facilita o uso de funcionalidades complexas com recursos nativos.
- **Backend (API)**: Desenvolvido em **Go (Golang)**, pela sua performance superior e menor custo de infraestrutura.
- **Comunicação**: Via **REST API**, com possibilidade futura de gRPC.
- **Armazenamento**: Centralizado na API.

## 5. Modo de Uso

- **Conectividade**: O aplicativo funcionará 100% online.
- **Autenticação**: Login obrigatório por e-mail e senha.
  - Implementado via **JWT**, pela facilidade e leveza para aplicações com baixa complexidade.
- **Persistência de Dados**: Todos os dados serão armazenados na nuvem (API).
- **Backup**: Realizado automaticamente no servidor, sem necessidade de ação do usuário.

## 6. Banco de Dados

- **SGBD**: PostgreSQL.
- **Justificativa**: Escolha de banco relacional por priorizar consistência e estrutura. A complexidade de um banco NoSQL não é necessária neste estágio do projeto.
- **Escalabilidade**: Inicialmente será usada **escala vertical**. A aplicação poderá ser migrada para soluções mais complexas, como replicação ou particionamento, apenas quando houver necessidade real.

## 7. Infraestrutura

- **Hospedagem**: O backend será hospedado em uma **VPS básica**, suficiente para o estágio inicial.
- **Motivação**: O Go apresenta excelente performance em ambientes com recursos limitados, reduzindo custos.

## 8. Monitoramento e Logs

- **Monitoramento**: 
  - Será utilizado **Prometheus** para coleta de métricas, e **Grafana** para visualização e alertas.
  - Justificativa: Ferramentas open-source que oferecem excelente flexibilidade e controle, alinhadas com o objetivo de manter baixo custo.

- **Logs**:
  - A aplicação usará bibliotecas como **Logrus**, e os eventos importantes serão expostos como métricas para o Prometheus.
  - Isso permite integração direta com Grafana para detecção de erros e eventos anormais.

## 9. Escalabilidade e Manutenção

- **Escalabilidade**:
  - Inicialmente, será adotada **escalabilidade vertical**, aumentando os recursos da VPS conforme o crescimento da demanda.
  - A base de dados PostgreSQL será suficiente, sem necessidade de replicação ou sharding no início.

- **CI/CD e Deploy**:
  - Será implementada uma esteira de **CI/CD** para automatizar testes e deploys.
  - O modelo de **deploy blue/green** será utilizado para garantir zero downtime.
  - Serão implementados **testes unitários e de integração** com ferramentas nativas do Go.

---

Este documento cobre a versão inicial do projeto e será expandido conforme novos requisitos e necessidades evoluírem.

