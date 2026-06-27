import { Button } from "../../components/ui/Button";
import { Shield } from "lucide-react";

interface PrivacyPolicyProps {
  onAccept: () => void;
}

export default function PrivacyPolicy({ onAccept }: PrivacyPolicyProps) {
  return (
    <div className="min-h-screen bg-gray-50 dark:bg-gray-900 py-12 px-4">
      <div className="max-w-3xl mx-auto">
        <div className="bg-white dark:bg-gray-800 rounded-xl border border-gray-200 dark:border-gray-700 shadow-sm">
          <div className="p-8 border-b border-gray-200 dark:border-gray-700">
            <h1 className="text-2xl font-bold text-gray-900 dark:text-gray-100">
              Política de Privacidade
            </h1>
            <p className="text-sm text-gray-500 dark:text-gray-400 mt-1">
              Última atualização: 19 de junho de 2026
            </p>
          </div>

          <div className="p-8 max-h-[60vh] overflow-y-auto prose prose-sm dark:prose-invert">
            <section className="mb-6">
              <h2 className="text-lg font-semibold text-gray-800 dark:text-gray-200 mb-2">
                1. Introdução
              </h2>
              <p className="text-gray-600 dark:text-gray-400">
                A equipe ENIAC Systems ("nós", "nosso") está comprometida em proteger sua privacidade. Esta Política de Privacidade explica como o <strong>FileENIAC</strong> coleta, usa, armazena e protege suas informações pessoais.
              </p>
            </section>

            <section className="mb-6">
              <h2 className="text-lg font-semibold text-gray-800 dark:text-gray-200 mb-2">
                2. Dados Coletados
              </h2>
              <p className="text-gray-600 dark:text-gray-400">
                O FileENIAC coleta os seguintes tipos de dados:
              </p>
              <ul className="list-disc pl-5 text-gray-600 dark:text-gray-400 space-y-1 mt-2">
                <li><strong>Dados de sessão:</strong> nome, descrição e caminho do workspace configurado por você.</li>
                <li><strong>Credenciais GitHub:</strong> tokens de acesso ao GitHub, que são encriptados localmente usando AES-256-GCM e armazenados apenas no seu dispositivo.</li>
                <li><strong>Dados de projeto:</strong> configurações de projetos, servidores FTP, histórico de deploys e sincronizações.</li>
                <li><strong>Dados de uso:</strong> eventos de health check e métricas de uso para diagnóstico.</li>
              </ul>
            </section>

            <section className="mb-6">
              <h2 className="text-lg font-semibold text-gray-800 dark:text-gray-200 mb-2">
                3. Armazenamento e Segurança
              </h2>
              <p className="text-gray-600 dark:text-gray-400">
                Todos os dados são armazenados localmente no seu dispositivo. Suas credenciais GitHub são encriptadas usando AES-256-GCM com uma chave de criptografia derivada usando Argon2id. A chave mestra nunca é armazenada em plaintext.
              </p>
              <p className="text-gray-600 dark:text-gray-400 mt-2">
                Não transmitimos suas credenciais ou dados pessoais para nossos servidores. Todas as operações de autenticação com o GitHub são realizadas diretamente entre o aplicativo e a API do GitHub.
              </p>
            </section>

            <section className="mb-6">
              <h2 className="text-lg font-semibold text-gray-800 dark:text-gray-200 mb-2">
                4. Uso dos Dados
              </h2>
              <p className="text-gray-600 dark:text-gray-400">
                Utilizamos seus dados apenas para:
              </p>
              <ul className="list-disc pl-5 text-gray-600 dark:text-gray-400 space-y-1 mt-2">
                <li>Gerenciar suas sessões e configurações de workspace;</li>
                <li>Realizar operações de deploy e sincronização conforme solicitado por você;</li>
                <li>Fornecer funcionalidades de histórico e análise;</li>
                <li>Diagnosticar problemas e melhorar a estabilidade do aplicativo.</li>
              </ul>
            </section>

            <section className="mb-6">
              <h2 className="text-lg font-semibold text-gray-800 dark:text-gray-200 mb-2">
                5. Compartilhamento de Dados
              </h2>
              <p className="text-gray-600 dark:text-gray-400">
                Não vendemos, alugamos ou compartilhamos suas informações pessoais com terceiros, exceto:
              </p>
              <ul className="list-disc pl-5 text-gray-600 dark:text-gray-400 space-y-1 mt-2">
                <li>Quando exigido por lei ou ordem judicial;</li>
                <li>Para proteger nossos direitos legais;</li>
                <li>Com sua autorização explícita.</li>
              </ul>
            </section>

            <section className="mb-6">
              <h2 className="text-lg font-semibold text-gray-800 dark:text-gray-200 mb-2">
                6. Retenção de Dados
              </h2>
              <p className="text-gray-600 dark:text-gray-400">
                Os dados são retidos enquanto você mantiver uma conta ativa no aplicativo. Você pode solicitar a exclusão de seus dados a qualquer momento entrando em contato pelo email <strong>contato@eniacsystems.com.br</strong>. A exclusão será realizada no prazo de 30 dias.
              </p>
            </section>

            <section className="mb-6">
              <h2 className="text-lg font-semibold text-gray-800 dark:text-gray-200 mb-2">
                7. Seus Direitos (LGPD - Lei Geral de Proteção de Dados)
              </h2>
              <p className="text-gray-600 dark:text-gray-400">
                Em conformidade com a Lei Geral de Proteção de Dados (Lei nº 13.709/2018), você tem direito a:
              </p>
              <ul className="list-disc pl-5 text-gray-600 dark:text-gray-400 space-y-1 mt-2">
                <li><strong>Acesso:</strong> solicitar acesso aos seus dados pessoais;</li>
                <li><strong>Correção:</strong> solicitar correção de dados incompletos ou desatualizados;</li>
                <li><strong>Exclusão:</strong> solicitar a exclusão dos seus dados;</li>
                <li><strong>Portabilidade:</strong> receber seus dados em formato estruturado;</li>
                <li><strong>Revogação:</strong> revogar consentimento a qualquer momento.</li>
              </ul>
            </section>

            <section className="mb-6">
              <h2 className="text-lg font-semibold text-gray-800 dark:text-gray-200 mb-2">
                8. Cookies e Tecnologias Similaress
              </h2>
              <p className="text-gray-600 dark:text-gray-400">
                O FileENIAC não utiliza cookies de rastreamento. Utilizamos apenas localStorage para armazenar configurações do usuário, como preferências de tema e dados de sessão.
              </p>
            </section>

            <section className="mb-6">
              <h2 className="text-lg font-semibold text-gray-800 dark:text-gray-200 mb-2">
                9. Alterações desta Política
              </h2>
              <p className="text-gray-600 dark:text-gray-400">
                Esta política pode ser atualizada periodicamente. Notificaremos sobre alterações significativas através do aplicativo. A data da última atualização está indicada no topo deste documento.
              </p>
            </section>

            <section>
              <h2 className="text-lg font-semibold text-gray-800 dark:text-gray-200 mb-2">
                10. Contato do DPO
              </h2>
              <p className="text-gray-600 dark:text-gray-400">
                Para questões sobre privacidade e proteção de dados, entre em contato com nosso Encarregado de Proteção de Dados (DPO):
              </p>
              <p className="text-gray-600 dark:text-gray-400 mt-1">
                <strong>Email:</strong> privacidade@eniacsystems.com.br
              </p>
            </section>
          </div>

          <div className="p-8 border-t border-gray-200 dark:border-gray-700 bg-gray-50 dark:bg-gray-800">
            <div className="flex gap-3 justify-end">
              <Button variant="primary" onClick={onAccept} icon={<Shield className="h-4 w-4" />}>
                Li e aceito a Política de Privacidade
              </Button>
            </div>
          </div>
        </div>
      </div>
    </div>
  );
}