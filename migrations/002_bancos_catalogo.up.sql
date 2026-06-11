CREATE TABLE bancos_catalogo (
    id     SERIAL PRIMARY KEY,
    nombre TEXT UNIQUE NOT NULL
);

INSERT INTO bancos_catalogo (nombre) VALUES
    ('Banco de Chile'),
    ('Banco Santander Chile'),
    ('Banco BCI'),
    ('BancoEstado'),
    ('Banco Itaú Chile'),
    ('Scotiabank Chile'),
    ('Banco BICE'),
    ('Banco Security'),
    ('Banco Consorcio'),
    ('Banco Internacional'),
    ('Banco Ripley'),
    ('Banco Falabella'),
    ('HSBC Bank Chile'),
    ('Banco BTG Pactual Chile'),
    ('JP Morgan Chase Bank N.A.'),
    ('Deutsche Bank AG'),
    ('Bank of China Agencia en Chile'),
    ('Coopeuch'),
    ('Tenpo Prepago'),
    ('Mach');
