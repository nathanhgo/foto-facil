---
name: worcap-prep
description: >-
  Ajuda a preparar e manter atualizado o material de apresentação do
  FotoFácil para a WORCAP (workshop do INPE), incluindo abstract, roteiro de
  demo e respostas para perguntas técnicas esperadas. Use quando o usuário
  pedir ajuda para preparar a apresentação, o pitch, o roteiro de demo, o
  abstract ou slides sobre o projeto para a WORCAP ou outro evento acadêmico.
---

# Preparar apresentação na WORCAP

Mantém [architecture-docs/worcap.md](../../../architecture-docs/worcap.md) atualizado e ajuda a gerar material derivado (slides, falas, respostas).

## Passo a passo

1. **Atualizar o estado do projeto**: releia `architecture-docs/mvp.md` e `architecture-docs/stack.md` para saber quais nós/features estão realmente prontos hoje — nunca inclua na demo algo que ainda não funciona.

2. **Atualizar `architecture-docs/worcap.md`**:
   - Ajuste o roteiro de demo para usar apenas nós já implementados (Fase 1 completa; adicione passos da Fase 2 conforme forem ficando prontos).
   - Mantenha a seção "Gaps conhecidos" honesta e atualizada — é melhor a pessoa que apresenta already saber a limitação do que ser pega de surpresa.

3. **Gerar material derivado quando solicitado**:
   - **Abstract/resumo**: use as duas versões (técnica e acessível) em `worcap.md` como base, ajustando tamanho/tom conforme o formato exigido pelo evento (pôster, comunicação oral, etc.).
   - **Slides**: estruture como Problema → Arquitetura (DAG, nós) → Demo ao vivo → Diferencial acadêmico (nós de PDI) → Roadmap. Use os diagramas de `architecture-docs/general_architecture.md` como referência visual.
   - **Perguntas e respostas**: use a seção "Gaps conhecidos" de `worcap.md` para preparar respostas honestas e diretas sobre limitações atuais.

4. **Não inventar resultados**: se o usuário pedir métricas de desempenho (ex.: tempo de processamento, PSNR) que ainda não foram medidas, sinalize isso em vez de estimar um número — sugira rodar a medição real primeiro.

## Notas

- Todo conteúdo gerado é em pt-br, exceto se o usuário pedir explicitamente uma versão em inglês (comum em resumos para eventos com público internacional).
- Ao final, se `worcap.md` foi alterado, registre em `architecture-docs/logs.md`.
