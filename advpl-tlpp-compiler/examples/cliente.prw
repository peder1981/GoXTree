//----------------------------------------------------------
// Exemplo de código AdvPL para testar o compilador
// Arquivo: cliente.prw
// Descrição: Classe para gerenciamento de clientes
//----------------------------------------------------------

#INCLUDE 'PROTHEUS.CH'
#INCLUDE 'TOTVS.CH'

/*/{Protheus.doc} Cliente
Classe para gerenciamento de clientes
@type class
@author Compilador AdvPL/TLPP
@since 2025-03-04
/*/
Class Cliente From LongNameClass

    // Atributos da classe
    Data cCodigo
    Data cNome
    Data cEndereco
    Data cCidade
    Data cUF
    Data cCEP
    Data cTelefone
    Data cEmail
    Data dCadastro
    Data nLimite
    Data lAtivo

    // Métodos da classe
    Method New(cCodigo, cNome) Constructor
    Method SetEndereco(cEndereco, cCidade, cUF, cCEP)
    Method SetContato(cTelefone, cEmail)
    Method SetLimite(nLimite)
    Method Ativar()
    Method Desativar()
    Method ToString()
    Method Gravar()

EndClass

/*/{Protheus.doc} Cliente::New
Método construtor da classe Cliente
@type method
@author Compilador AdvPL/TLPP
@since 2025-03-04
@param cCodigo, Caractere, Código do cliente
@param cNome, Caractere, Nome do cliente
@return Self, Objeto, Instância da classe Cliente
/*/
Method New(cCodigo, cNome) Class Cliente
    // Inicialização dos atributos
    ::cCodigo   := cCodigo
    ::cNome     := cNome
    ::dCadastro := Date()
    ::lAtivo    := .T.
    
    Return Self

/*/{Protheus.doc} Cliente::SetEndereco
Define o endereço do cliente
@type method
@author Compilador AdvPL/TLPP
@since 2025-03-04
@param cEndereco, Caractere, Endereço do cliente
@param cCidade, Caractere, Cidade do cliente
@param cUF, Caractere, UF do cliente
@param cCEP, Caractere, CEP do cliente
@return Self, Objeto, Instância da classe Cliente
/*/
Method SetEndereco(cEndereco, cCidade, cUF, cCEP) Class Cliente
    ::cEndereco := cEndereco
    ::cCidade   := cCidade
    ::cUF       := cUF
    ::cCEP      := cCEP
    
    Return Self

/*/{Protheus.doc} Cliente::SetContato
Define as informações de contato do cliente
@type method
@author Compilador AdvPL/TLPP
@since 2025-03-04
@param cTelefone, Caractere, Telefone do cliente
@param cEmail, Caractere, Email do cliente
@return Self, Objeto, Instância da classe Cliente
/*/
Method SetContato(cTelefone, cEmail) Class Cliente
    ::cTelefone := cTelefone
    ::cEmail    := cEmail
    
    Return Self

/*/{Protheus.doc} Cliente::SetLimite
Define o limite de crédito do cliente
@type method
@author Compilador AdvPL/TLPP
@since 2025-03-04
@param nLimite, Numérico, Limite de crédito do cliente
@return Self, Objeto, Instância da classe Cliente
/*/
Method SetLimite(nLimite) Class Cliente
    ::nLimite := nLimite
    
    Return Self

/*/{Protheus.doc} Cliente::Ativar
Ativa o cliente
@type method
@author Compilador AdvPL/TLPP
@since 2025-03-04
@return Self, Objeto, Instância da classe Cliente
/*/
Method Ativar() Class Cliente
    ::lAtivo := .T.
    
    Return Self

/*/{Protheus.doc} Cliente::Desativar
Desativa o cliente
@type method
@author Compilador AdvPL/TLPP
@since 2025-03-04
@return Self, Objeto, Instância da classe Cliente
/*/
Method Desativar() Class Cliente
    ::lAtivo := .F.
    
    Return Self

/*/{Protheus.doc} Cliente::ToString
Retorna uma string com as informações do cliente
@type method
@author Compilador AdvPL/TLPP
@since 2025-03-04
@return cString, Caractere, String com as informações do cliente
/*/
Method ToString() Class Cliente
    Local cString := ""
    
    cString += "Código: " + ::cCodigo + CRLF
    cString += "Nome: " + ::cNome + CRLF
    
    If !Empty(::cEndereco)
        cString += "Endereço: " + ::cEndereco + CRLF
        cString += "Cidade: " + ::cCidade + " - " + ::cUF + CRLF
        cString += "CEP: " + ::cCEP + CRLF
    EndIf
    
    If !Empty(::cTelefone)
        cString += "Telefone: " + ::cTelefone + CRLF
    EndIf
    
    If !Empty(::cEmail)
        cString += "Email: " + ::cEmail + CRLF
    EndIf
    
    cString += "Data de Cadastro: " + DTOC(::dCadastro) + CRLF
    cString += "Limite de Crédito: " + Transform(::nLimite, "@E 999,999,999.99") + CRLF
    cString += "Status: " + IIF(::lAtivo, "Ativo", "Inativo") + CRLF
    
    Return cString

/*/{Protheus.doc} Cliente::Gravar
Grava os dados do cliente no banco de dados
@type method
@author Compilador AdvPL/TLPP
@since 2025-03-04
@return lSucesso, Lógico, Indica se a gravação foi bem-sucedida
/*/
Method Gravar() Class Cliente
    Local lSucesso := .F.
    Local aArea    := GetArea()
    
    // Abre a tabela de clientes
    DbSelectArea("SA1")
    DbSetOrder(1)
    
    // Verifica se o cliente já existe
    If DbSeek(xFilial("SA1") + ::cCodigo)
        // Atualiza o cliente
        RecLock("SA1", .F.)
    Else
        // Insere um novo cliente
        RecLock("SA1", .T.)
        SA1->A1_FILIAL := xFilial("SA1")
        SA1->A1_COD    := ::cCodigo
    EndIf
    
    // Atualiza os campos
    SA1->A1_NOME    := ::cNome
    SA1->A1_END     := ::cEndereco
    SA1->A1_MUN     := ::cCidade
    SA1->A1_EST     := ::cUF
    SA1->A1_CEP     := ::cCEP
    SA1->A1_TEL     := ::cTelefone
    SA1->A1_EMAIL   := ::cEmail
    SA1->A1_DTCAD   := ::dCadastro
    SA1->A1_LC      := ::nLimite
    SA1->A1_MSBLQL  := IIF(::lAtivo, "2", "1")
    
    // Confirma a gravação
    MsUnlock()
    lSucesso := .T.
    
    // Restaura a área
    RestArea(aArea)
    
    Return lSucesso

/*/{Protheus.doc} CriarCliente
Função para criar um novo cliente
@type function
@author Compilador AdvPL/TLPP
@since 2025-03-04
@param cCodigo, Caractere, Código do cliente
@param cNome, Caractere, Nome do cliente
@param cEndereco, Caractere, Endereço do cliente
@param cCidade, Caractere, Cidade do cliente
@param cUF, Caractere, UF do cliente
@param cCEP, Caractere, CEP do cliente
@param cTelefone, Caractere, Telefone do cliente
@param cEmail, Caractere, Email do cliente
@param nLimite, Numérico, Limite de crédito do cliente
@return oCliente, Objeto, Instância da classe Cliente
/*/
Function CriarCliente(cCodigo, cNome, cEndereco, cCidade, cUF, cCEP, cTelefone, cEmail, nLimite)
    Local oCliente := Cliente():New(cCodigo, cNome)
    
    // Define os dados do cliente
    oCliente:SetEndereco(cEndereco, cCidade, cUF, cCEP)
    oCliente:SetContato(cTelefone, cEmail)
    oCliente:SetLimite(nLimite)
    
    // Grava os dados
    If oCliente:Gravar()
        ConOut("Cliente " + cCodigo + " criado com sucesso!")
    Else
        ConOut("Erro ao criar o cliente " + cCodigo)
    EndIf
    
    Return oCliente

/*/{Protheus.doc} TestarCliente
Função para testar a classe Cliente
@type function
@author Compilador AdvPL/TLPP
@since 2025-03-04
@return Nil
/*/
Function TestarCliente()
    Local oCliente
    Local i
    Local aClientes := {}
    
    // Cria alguns clientes de teste
    For i := 1 To 5
        oCliente := Cliente():New("C" + StrZero(i, 5), "Cliente de Teste " + StrZero(i, 2))
        oCliente:SetEndereco("Rua Teste, " + StrZero(i*100, 4), "São Paulo", "SP", "01000-000")
        oCliente:SetContato("(11) 9999-9999", "cliente" + StrZero(i, 2) + "@teste.com.br")
        oCliente:SetLimite(i * 1000)
        
        If i % 2 == 0
            oCliente:Desativar()
        EndIf
        
        AAdd(aClientes, oCliente)
    Next i
    
    // Exibe os dados dos clientes
    For i := 1 To Len(aClientes)
        ConOut("------------------------------------------")
        ConOut("Cliente " + StrZero(i, 2))
        ConOut("------------------------------------------")
        ConOut(aClientes[i]:ToString())
    Next i
    
    Return Nil

/*/{Protheus.doc} u_TesteCliente
Função de usuário para testar a classe Cliente
@type user function
@author Compilador AdvPL/TLPP
@since 2025-03-04
@return Nil
/*/
User Function TesteCliente()
    Local oCliente := CriarCliente("C00001", "CLIENTE TESTE", "RUA TESTE, 123", "SÃO PAULO", "SP", "01000-000", "(11) 9999-9999", "teste@teste.com.br", 10000)
    
    // Exibe os dados do cliente
    MsgInfo(oCliente:ToString(), "Dados do Cliente")
    
    Return Nil
