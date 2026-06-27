import { Button } from "../../components/ui/Button";
import { CheckCircle } from "lucide-react";

interface TermsProps {
  onAccept: () => void;
}

export default function TermsOfUse({ onAccept }: TermsProps) {
  return (
    <div className="min-h-screen bg-gray-50 dark:bg-gray-900 py-12 px-4">
      <div className="max-w-3xl mx-auto">
        <div className="bg-white dark:bg-gray-800 rounded-xl border border-gray-200 dark:border-gray-700 shadow-sm">
          <div className="p-8 border-b border-gray-200 dark:border-gray-700">
            <h1 className="text-2xl font-bold text-gray-900 dark:text-gray-100">
              Termos de Uso
            </h1>
            <p className="text-sm text-gray-500 dark:text-gray-400 mt-1">
              Última atualização: 19 de junho de 2026
            </p>
          </div>

          <div className="p-8 max-h-[60vh] overflow-y-auto prose prose-sm dark:prose-invert">
            <section className="mb-6">
              <h2 className="text-lg font-semibold text-gray-800 dark:text-gray-200 mb-2">
                1. Aceitação dos Termos
              </h2>
              <p className="text-gray-600 dark:text-gray-400">
                Ao acessar e utilizar o <strong>FileENIAC</strong>, você concorda com os presentes Termos de Uso. Se você não concorda com algum dos termos, não utilize o aplicativo.
              </p>
            </section>

            <section className="mb-6">
              <h2 className="text-lg font-semibold text-gray-800 dark:text-gray-200 mb-2">
                2. Descrição do Serviço
              </h2>
              <p className="text-gray-600 dark:text-gray-400">
                O FileENIAC é uma ferramenta de gerenciamento de projetos, deploy FTP e sincronização de workspaces desenvolvida pela equipe ENIAC Systems. O serviço permite gerenciar múltiplos projetos, servidores FTP, deploys e históricos de alterações.
              </p>
            </section>

            <section className="mb-6">
              <h2 className="text-lg font-semibold text-gray-800 dark:text-gray-200 mb-2">
                3. Conta e Segurança
              </h2>
              <p className="text-gray-600 dark:text-gray-400">
                Você é responsável por manter a confidencialidade de suas credenciais de acesso ao GitHub e outras informações de configuração. Todas as operações realizadas através do FileENIAC são de sua inteira responsabilidade.
              </p>
            </section>

            <section className="mb-6">
              <h2 className="text-lg font-semibold text-gray-800 dark:text-gray-200 mb-2">
                4. Dados e Privacidade
              </h2>
              <p className="text-gray-600 dark:text-gray-400">
                O FileENIAC armazena dados localmente no seu dispositivo, incluindo configurações de sessão, histórico de deploys e credenciais encriptadas. Para mais detalhes, consulte nossa Política de Privacidade.
              </p>
            </section>

            <section className="mb-6">
              <h2 className="text-lg font-semibold text-gray-800 dark:text-gray-200 mb-2">
                5. Uso Acceptável
              </h2>
              <p className="text-gray-600 dark:text-gray-400">
                Você se compromete a utilizar o FileENIAC apenas para fins legais e de acordo com os presentes termos. É proibido utilizar o serviço para atividades ilegais, violação de direitos de terceiros, ou qualquer uso que possa prejudicar o funcionamento do aplicativo ou de terceiros.
              </p>
            </section>

            <section className="mb-6">
              <h2 className="text-lg font-semibold text-gray-800 dark:text-gray-200 mb-2">
                6. Isenção de Garantia
              </h2>
              <p className="text-gray-600 dark:text-gray-400">
                O FileENIAC é fornecido "como está", sem garantias de qualquer natureza, expressas ou implícitas. Não garantimos que o serviço será ininterrupto, seguro ou livre de erros.
              </p>
            </section>

            <section className="mb-6">
              <h2 className="text-lg font-semibold text-gray-800 dark:text-gray-200 mb-2">
                7. Limitação de Responsabilidade
              </h2>
              <p className="text-gray-600 dark:text-gray-400">
                Em nenhuma circunstância a equipe ENIAC Systems será responsável por quaisquer danos diretos, indiretos, incidentais, especiais ou consequentes decorrentes do uso ou da incapacidade de uso do FileENIAC.
              </p>
            </section>

            <section className="mb-6">
              <h2 className="text-lg font-semibold text-gray-800 dark:text-gray-200 mb-2">
                8. Alterações dos Termos
              </h2>
              <p className="text-gray-600 dark:text-gray-400">
                Reservamo-nos o direito de modificar estes termos a qualquer momento. As alterações entrarão em vigor imediatamente após a publicação. O uso continuado do serviço após as alterações constitui aceitação dos novos termos.
              </p>
            </section>

            <section>
              <h2 className="text-lg font-semibold text-gray-800 dark:text-gray-200 mb-2">
                9. Contato
              </h2>
              <p className="text-gray-600 dark:text-gray-400">
                Para dúvidas sobre estes termos, entre em contato pelo email: <strong>contato@eniacsystems.com.br</strong>
              </p>
            </section>
          </div>

          <div className="p-8 border-t border-gray-200 dark:border-gray-700 bg-gray-50 dark:bg-gray-800">
            <div className="flex gap-3 justify-end">
              <Button variant="primary" onClick={onAccept} icon={<CheckCircle className="h-4 w-4" />}>
                Li e aceito os Termos de Uso
              </Button>
            </div>
          </div>
        </div>
      </div>
    </div>
  );
}