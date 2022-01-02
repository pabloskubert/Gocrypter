# Gocrypter
Crypter em golang (POC)

## Estágios do crypter

1. Comprimi o arquivo malicioso usando a ZLIB
2. Criptografa os bytes resultantes do processo de compressão usando AES-GCM com chave 256-bit long
3. Converte o resultado da criptografia byte-a-byte para uint8 e grava no stub_template.go na seção "PAYLOAD"
4. Gera um arquivo main.go com base no stub_template.go no diretório "stub"
5. Compila o arquivo com o stub + payload

## Estágios da infecção

1. Descriptografa os bytes da carga que foram armazenados em uma variável no stub (i.e payload)
2. Descomprimi os bytes do payload usando a ZLIB
3. Grava a carga maliciosa no disco rígido para o diretório temporário do sistema operacional
4. Executa o arquivo que foi gravado


## Efetividade

A efetividade do crypter é consistente até que a carga maliciosa seja gravado no disco rígido, somente 
esta técnica sozinha não é capaz de ser efetiva contra EDR's e AV's inteligentes, é necessário combiná-la 
com outras de injenção, como `PROCESS HOLLOWING` ou qualquer outra. 
