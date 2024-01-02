package types

type BusinessActivity string

const (
	BusinessActivityA BusinessActivity = "A" //AGRICULTURA, PECUÁRIA, PROD. FLORESTAL, PESCA
	BusinessActivityB BusinessActivity = "B" //INDÚSTRIAS EXTRATIVAS
	BusinessActivityC BusinessActivity = "C" //INDÚSTRIAS DE TRANSFORMAÇÃO
	BusinessActivityD BusinessActivity = "D" //ELETRICIDADE E GÁS
	BusinessActivityE BusinessActivity = "E" //ÁGUA, ESGOTO, ATIVIDADES DE GESTÃO DE RESÍDUOS
	BusinessActivityF BusinessActivity = "F" //CONSTRUÇÃO
	BusinessActivityG BusinessActivity = "G" //COMÉRCIO, REPARAÇÃO DE VEÍCULOS
	BusinessActivityH BusinessActivity = "H" //TRANSPORTE, ARMAZENAGEM E CORREIO
	BusinessActivityI BusinessActivity = "I" //ALOJAMENTO E ALIMENTAÇÃO
	BusinessActivityJ BusinessActivity = "J" //INFORMAÇÃO E COMUNICAÇÃO
	BusinessActivityK BusinessActivity = "K" //ATIVIDADES FINANCEIRAS, DE SEGUROS E SERVIÇOS
	BusinessActivityL BusinessActivity = "L" //ATIVIDADES IMOBILIÁRIAS
	BusinessActivityM BusinessActivity = "M" //ATIVIDADES PROFISSIONAIS,CIENTÍFICAS E TÉCN.
	BusinessActivityN BusinessActivity = "N" //ATIVIDADES ADMINISTRATIVAS E SERVIÇOS COMPL.
	BusinessActivityO BusinessActivity = "O" //ADMIN. PÚBLICA, DEFESA E SEGURIDADE SOCIAL
	BusinessActivityP BusinessActivity = "P" //EDUCAÇÃO
	BusinessActivityQ BusinessActivity = "Q" //SAÚDE HUMANA E SERVIÇOS SOCIAIS
	BusinessActivityR BusinessActivity = "R" //ARTES, CULTURA, ESPORTE E RECREAÇÃO
	BusinessActivityS BusinessActivity = "S" //OUTRAS ATIVIDADES DE SERVIÇOS
	BusinessActivityT BusinessActivity = "T" //SERVIÇOS DOMÉSTICOS
	BusinessActivityU BusinessActivity = "U" //ORGANISMOS INTERNACIONAIS
)