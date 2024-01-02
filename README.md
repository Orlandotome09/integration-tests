This is a README file


## 1.Como atualizar mocks (após alteração de alguma interface da pasta core/_interfaces)

```
mockery --dir src/core/_interfaces/ --all --output src/core/_mocks 
```

## 2.Acessando o Swagger local
```
http://localhost:{porta}/compliance-int/swagger/index.html
```

## 3.Rodando a aplicação localmente

Caso não tenha instalado as dependências, uso o seguinte comando na **raiz** do projeto:
```
go mod download
```

### 3.1.Subir as dependências com ***docker-compose***
O seguinte comando levarantará os serviços externos que a aplicação usa, como bancos de dados, filas e serviços mockados (usar na **raiz** do projeto):

```
sudo sh compose.sh
```
Se houver algum problema com ***conflito de portas***, é preciso derrubar os conteineres em execução. Para derrubar todos os conteineres:
```
sudo docker stop $(docker ps -q)
```

### 3.2.Carregar as variáveis de ambiente
O seguinte arquivo possui as variáveis de ambientes para rodar a aplicação localmente:

```
arquivo .env
```

Para facilitar o uso de variáveis de ambiente no Goland:
- Instalar o plugin ***EnvFile*** de ***Boris Pierov***

### 3.3.Rodar a app
I. Clicar em ***Select Run and Debug Configuration***

II.Selecione ***"+"*** e ***Go Build***

III.Em ***Run King*** selecione ***Directory***

IV.Em ***Directory*** selecione o caminho da pasta ***src*** da aplicação

V.Na aba ***EnvFile*** (precisa ter instalado o plugin do passo 3.2), selecione ***Enable EnvFile*** e, clique em ***"+"***

VI.Selecione o arquivo ***.env*** da raiz do projeto (pode ser que o arquivo esteja escondido, neste caso, clique em ***"show hidden files"***)

## 4.Rodando os Testes Unitários
Na raiz do projeto, rodar o comando:
```
go test ./...
```
Se algum teste quebrar, lembre-se de fazer a correção antes de abrir qualquer ***Pull Request***.

##  5.Testes de Integração
Rodar o comando na raiz do projeto:
```
sudo sh integration_test.sh
```

Se algum teste falhar, os conteineres continuarão de pé para que, caso seja necessário, seja possível fazer o debug dentro dos mesmos.

Para ***debugar*** um contêiner:

I.Pegue o ***id*** do contêiner com ```sudo docker ps```

II.Veja os logs com ```sudo docker logs {id}```


## 6.Testes de Integração com Karate
Durante o desenvolvimento, é possível que o desenvolvedor queira verificar se os testes de integração estão funcionando, sem usar o scrip do ```integration_test.sh``` (que roda todos os testes e pode demorar).

Para isso:

I.Subir as depedências da app (sh compose.sh) e rodar a app.

II.No Visual Studio Code, abrir o diretório onde está o arquivo ***karate-config.js***:
```
{raiz}/tests/integration
```
III.No Visual Studio Code, ir em:
```
Preferences -> Settings -> Workspace
```
IV.Procurar por:
```
Karate Runner -> Karate Jar: Command Line Args
```

V.Na linha digitável, abaixo do ***"command line args"***, apontar os jars que estão na pasta:
```
{raiz}/tests/karate
```

São dois jars: ***karate.jar*** e ***KarateUtils.jar***.

VI.Instalar a extensão ***Karate Runner***, de ***Kirk Slota***.

VII.Abrir o arquivo ***.feature*** que deseja executar e clicar em ***Karate:Debug***

