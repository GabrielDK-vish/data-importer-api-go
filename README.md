# Desafio Técnico — Full Stack (Golang)

##  Desafio proposto:
> Você deverá criar um **importador para uma base de dados (Postgres)** que deverá ser feito em **Golang** para armazenar os dados do arquivo enviado.  
> Será avaliado:  
> - **Normalização dos dados** na base de dados  
> - **Performance** do importador  
>  
> Além disso:  
> - Criar uma **API em Golang** com:
>   - Endpoint de **autenticação**
>   - Endpoints de **consulta dos dados importados**  
> - **Diferencial:** Criar um **frontend em React** mostrando indicadores totalizadores, agrupamentos por categorias, recursos, clientes e meses de cobrança.  
> - Necessário **publicar em algum link** para avaliação, com documentação de execução.

---

##  Solução proposta

1. **Input (CLI em Go)**
   - Leitura do arquivo.
   - Conversão e normalização de dados.
   - Envio em lotes para Postgres via `pgx.CopyFrom`.

2. **Banco de Dados (Postgres)**
   Estrutura normalizada em 4 entidades principais:
   - **partners**: dados do parceiro (PartnerId, PartnerName, MpnId, Tier2MpnId).
   - **customers**: dados do cliente (CustomerId, CustomerName, CustomerDomainName, Country).
   - **products**: catálogo de serviços/recursos (ProductId, SkuId, SkuName, ProductName, MeterType, Category, SubCategory, UnitType).
   - **usages**: registros de consumo e faturamento, vinculando `partner_id`, `customer_id` e `product_id`  
     (InvoiceNumber, ChargeStartDate, UsageDate, Quantity, UnitPrice, BillingPreTaxTotal, ResourceLocation, Tags, BenefitType).

   ➝ Essa separação garante **normalização** e facilita consultas.

3. **API (Go)**
   - Framework: `chi`.
   - Autenticação via **JWT**.
   - Endpoints:
     - `POST /auth/login` → autenticação
     - `GET /customers` → listar clientes
     - `GET /customers/{id}/usage` → consumo detalhado do cliente
     - `GET /reports/billing/monthly` → total por mês
     - `GET /reports/billing/by-product` → agrupado por produto/serviço
     - `GET /reports/billing/by-partner` → agrupado por parceiro

4. **Frontend (React) **
   - Dashboard com indicadores:
     - Faturamento total por mês
     - Ranking de clientes por consumo
     - Distribuição por produtos/recursos
   - Gráficos via Recharts.

5. **Infraestrutura**
   - `docker-compose` para Postgres + API + Front.
   - Deploy em **Render/Railway** (API) e **Vercel** (Frontend).
   - Migrations via `golang-migrate`.


---

## Estrutura inicial do projeto

billing-importer-go/<br>
├── cmd/<br>
│ ├── importer/ # CLI de importação<br>
│ └── api/ # Servidor HTTP<br>
├── internal/<br>
│ ├── api/ # Handlers e rotas<br>
│ ├── db/ # Conexão e queries<br>
│ ├── models/ # Structs e DTOs<br>
│ └── auth/ # JWT e segurança<br>
├── frontend/ # (opcional) React App<br>
├── migrations/ # Scripts SQL<br>
├── docker-compose.yml<br>
├── Dockerfile<br>
└── README.md<br>


---

## 🌐 Publicação
- **API:** _(N/A)_  
- **Frontend:** _(N/A)_  

---

## 📝 Documentação
- [API (endpoints e exemplos)](./docs/api.md)  
- [Guia do Importador](./docs/importer.md)  
- [Migrations](./migrations)  
