// SPDX-License-Identifier: MIT
import { useState } from "react";
import { useNavigate } from "react-router-dom";
import { Button } from "../../components/ui/Button";
import { Shield, FileText, Lock } from "lucide-react";
import { STORAGE_KEYS, storageSet } from "../../api/storage";

interface LGPDConsentProps {
  onComplete: () => void;
}

export default function LGPDConsent({ onComplete }: LGPDConsentProps) {
  const navigate = useNavigate();
  const [step, setStep] = useState<"intro" | "terms" | "privacy">("intro");
  const [termsAccepted, setTermsAccepted] = useState(false);
  const [privacyAccepted, setPrivacyAccepted] = useState(false);
  const [dataProcessingAccepted, setDataProcessingAccepted] = useState(false);

  const canProceed = termsAccepted && privacyAccepted && dataProcessingAccepted;

  function handleComplete() {
    storageSet(STORAGE_KEYS.lgpdConsent, JSON.stringify({
      agreed: true,
      termsAccepted: true,
      privacyAccepted: true,
      dataProcessingAccepted: true,
      consentedAt: new Date().toISOString(),
    }));
    onComplete();
  }

  if (step === "terms") {
    return (
      <div className="min-h-screen bg-gray-50 dark:bg-gray-900 py-8 px-4">
        <div className="max-w-3xl mx-auto">
          <div className="bg-white dark:bg-gray-800 rounded-xl border border-gray-200 dark:border-gray-700 shadow-sm">
            <div className="p-6 border-b border-gray-200 dark:border-gray-700">
              <h1 className="text-xl font-bold text-gray-900 dark:text-gray-100">
                Termos de Uso
              </h1>
              <p className="text-sm text-gray-500 dark:text-gray-400 mt-1">
                Por favor, leia attentamente antes de continuar.
              </p>
            </div>
            <div className="p-6 max-h-[50vh] overflow-y-auto">
              <div className="prose prose-sm dark:prose-invert">
                <section className="mb-4">
                  <h3 className="text-base font-semibold text-gray-800 dark:text-gray-200">1. Aceitação dos Termos</h3>
                  <p className="text-gray-600 dark:text-gray-400 text-sm">
                    Ao acessar e utilizar o FileENIAC, você concorda com os presentes Termos de Uso. Se você não concorda com algum dos termos, não utilize o aplicativo.
                  </p>
                </section>
                <section className="mb-4">
                  <h3 className="text-base font-semibold text-gray-800 dark:text-gray-200">2. Descrição do Serviço</h3>
                  <p className="text-gray-600 dark:text-gray-400 text-sm">
                    O FileENIAC é uma ferramenta de gerenciamento de projetos, deploy FTP e sincronização de workspaces.
                  </p>
                </section>
                <section className="mb-4">
                  <h3 className="text-base font-semibold text-gray-800 dark:text-gray-200">3. Conta e Segurança</h3>
                  <p className="text-gray-600 dark:text-gray-400 text-sm">
                    Você é responsável por manter a confidencialidade de suas credenciais de acesso ao GitHub.
                  </p>
                </section>
                <section className="mb-4">
                  <h3 className="text-base font-semibold text-gray-800 dark:text-gray-200">4. Dados e Privacidade</h3>
                  <p className="text-gray-600 dark:text-gray-400 text-sm">
                    O FileENIAC armazena dados localmente no seu dispositivo. Credenciais GitHub são encriptadas usando AES-256-GCM.
                  </p>
                </section>
                <section className="mb-4">
                  <h3 className="text-base font-semibold text-gray-800 dark:text-gray-200">5. Uso Acceptável</h3>
                  <p className="text-gray-600 dark:text-gray-400 text-sm">
                    É proibido utilizar o serviço para atividades ilegais ou que violem direitos de terceiros.
                  </p>
                </section>
                <section className="mb-4">
                  <h3 className="text-base font-semibold text-gray-800 dark:text-gray-200">6. Isenção de Garantia</h3>
                  <p className="text-gray-600 dark:text-gray-400 text-sm">
                    O FileENIAC é fornecido "como está", sem garantias de qualquer natureza.
                  </p>
                </section>
                <section>
                  <h3 className="text-base font-semibold text-gray-800 dark:text-gray-200">7. Contato</h3>
                  <p className="text-gray-600 dark:text-gray-400 text-sm">
                    Dúvidas: contato@eniacsystems.com.br
                  </p>
                </section>
              </div>
            </div>
            <div className="p-6 border-t border-gray-200 dark:border-gray-700 bg-gray-50 dark:bg-gray-800">
              <div className="flex items-center gap-3 mb-4">
                <input
                  type="checkbox"
                  id="termsAccepted"
                  checked={termsAccepted}
                  onChange={(e) => setTermsAccepted(e.target.checked)}
                  className="w-4 h-4 rounded border-gray-300 text-eniac-600 focus:ring-eniac-500"
                />
                <label htmlFor="termsAccepted" className="text-sm text-gray-700 dark:text-gray-300">
                  Li e aceito os Termos de Uso
                </label>
              </div>
              <div className="flex gap-3 justify-end">
                <Button variant="secondary" onClick={() => setStep("intro")}>
                  Voltar
                </Button>
                <Button variant="primary" onClick={() => setStep("privacy")} disabled={!termsAccepted}>
                  Continuar para Privacidade
                </Button>
              </div>
            </div>
          </div>
        </div>
      </div>
    );
  }

  if (step === "privacy") {
    return (
      <div className="min-h-screen bg-gray-50 dark:bg-gray-900 py-8 px-4">
        <div className="max-w-3xl mx-auto">
          <div className="bg-white dark:bg-gray-800 rounded-xl border border-gray-200 dark:border-gray-700 shadow-sm">
            <div className="p-6 border-b border-gray-200 dark:border-gray-700">
              <h1 className="text-xl font-bold text-gray-900 dark:text-gray-100">
                Política de Privacidade
              </h1>
              <p className="text-sm text-gray-500 dark:text-gray-400 mt-1">
                Suas informações e como as protegemos.
              </p>
            </div>
            <div className="p-6 max-h-[50vh] overflow-y-auto">
              <div className="prose prose-sm dark:prose-invert">
                <section className="mb-4">
                  <h3 className="text-base font-semibold text-gray-800 dark:text-gray-200">1. Dados Coletados</h3>
                  <p className="text-gray-600 dark:text-gray-400 text-sm">
                    Coletamos: dados de sessão (nome, workspace), credenciais GitHub (encriptadas), dados de projetos e histórico de deploys.
                  </p>
                </section>
                <section className="mb-4">
                  <h3 className="text-base font-semibold text-gray-800 dark:text-gray-200">2. Armazenamento e Segurança</h3>
                  <p className="text-gray-600 dark:text-gray-400 text-sm">
                    Todos os dados são armazenados localmente no seu dispositivo. Credenciais são encriptadas com AES-256-GCM + Argon2id. Não transmitimos credenciais para nossos servidores.
                  </p>
                </section>
                <section className="mb-4">
                  <h3 className="text-base font-semibold text-gray-800 dark:text-gray-200">3. Seus Direitos (LGPD)</h3>
                  <p className="text-gray-600 dark:text-gray-400 text-sm">
                    Você tem direito a: acesso, correção, exclusão, portabilidade e revogação do consentimento. Contato: privacidade@eniacsystems.com.br
                  </p>
                </section>
                <section className="mb-4">
                  <h3 className="text-base font-semibold text-gray-800 dark:text-gray-200">4. Retenção</h3>
                  <p className="text-gray-600 dark:text-gray-400 text-sm">
                    Dados são retidos enquanto a conta estiver ativa. Exclusão pode ser solicitada a qualquer momento pelo email privacidade@eniacsystems.com.br.
                  </p>
                </section>
                <section>
                  <h3 className="text-base font-semibold text-gray-800 dark:text-gray-200">5. Cookies</h3>
                  <p className="text-gray-600 dark:text-gray-400 text-sm">
                    Não utilizamos cookies de rastreamento. Apenas localStorage para preferências do usuário.
                  </p>
                </section>
              </div>
            </div>
            <div className="p-6 border-t border-gray-200 dark:border-gray-700 bg-gray-50 dark:bg-gray-800">
              <div className="flex items-center gap-3 mb-4">
                <input
                  type="checkbox"
                  id="privacyAccepted"
                  checked={privacyAccepted}
                  onChange={(e) => setPrivacyAccepted(e.target.checked)}
                  className="w-4 h-4 rounded border-gray-300 text-eniac-600 focus:ring-eniac-500"
                />
                <label htmlFor="privacyAccepted" className="text-sm text-gray-700 dark:text-gray-300">
                  Li e aceito a Política de Privacidade
                </label>
              </div>
              <div className="flex items-center gap-3 mb-4">
                <input
                  type="checkbox"
                  id="dataProcessingAccepted"
                  checked={dataProcessingAccepted}
                  onChange={(e) => setDataProcessingAccepted(e.target.checked)}
                  className="w-4 h-4 rounded border-gray-300 text-eniac-600 focus:ring-eniac-500"
                />
                <label htmlFor="dataProcessingAccepted" className="text-sm text-gray-700 dark:text-gray-300">
                  Consinto com o armazenamento e processamento dos meus dados conforme descrito na Política de Privacidade
                </label>
              </div>
              <div className="flex gap-3 justify-end">
                <Button variant="secondary" onClick={() => setStep("terms")}>
                  Voltar
                </Button>
                <Button variant="primary" onClick={handleComplete} disabled={!canProceed}>
                  Concluir
                </Button>
              </div>
            </div>
          </div>
        </div>
      </div>
    );
  }

  return (
    <div className="min-h-screen bg-gradient-to-br from-eniac-900 to-eniac-700 flex items-center justify-center p-4">
      <div className="max-w-lg w-full bg-white dark:bg-gray-800 rounded-2xl shadow-2xl p-8">
        <div className="text-center mb-8">
          <div className="w-16 h-16 bg-eniac-100 dark:bg-eniac-900 rounded-2xl flex items-center justify-center mx-auto mb-4">
            <Shield className="h-8 w-8 text-eniac-600 dark:text-eniac-400" />
          </div>
          <h1 className="text-2xl font-bold text-gray-900 dark:text-gray-100">
            Bienvenido ao FileENIAC
          </h1>
          <p className="text-gray-500 dark:text-gray-400 mt-2 text-sm">
            Antes de continuar, é necessário aceitar nossos Termos de Uso e Política de Privacidade, em conformidade com a LGPD.
          </p>
        </div>

        <div className="space-y-3 mb-8">
          <button
            onClick={() => setStep("terms")}
            className="w-full flex items-center gap-3 p-4 bg-gray-50 dark:bg-gray-700 rounded-xl hover:bg-gray-100 dark:hover:bg-gray-600 transition-colors text-left"
          >
            <FileText className="h-5 w-5 text-eniac-600 dark:text-eniac-400" />
            <div>
              <p className="font-medium text-gray-900 dark:text-gray-100">Termos de Uso</p>
              <p className="text-xs text-gray-500 dark:text-gray-400">Regras e condições de uso do aplicativo</p>
            </div>
          </button>

          <button
            onClick={() => setStep("privacy")}
            className="w-full flex items-center gap-3 p-4 bg-gray-50 dark:bg-gray-700 rounded-xl hover:bg-gray-100 dark:hover:bg-gray-600 transition-colors text-left"
          >
            <Lock className="h-5 w-5 text-eniac-600 dark:text-eniac-400" />
            <div>
              <p className="font-medium text-gray-900 dark:text-gray-100">Política de Privacidade</p>
              <p className="text-xs text-gray-500 dark:text-gray-400">Como protegemos seus dados (LGPD)</p>
            </div>
          </button>
        </div>

        <Button
          variant="primary"
          size="lg"
          className="w-full"
          onClick={() => setStep("terms")}
          icon={<Shield className="h-5 w-5" />}
        >
          Ler e Aceitar Termos
        </Button>

        <div className="mt-6 pt-6 border-t border-gray-200 dark:border-gray-700">
          <button
            onClick={() => navigate("/")}
            className="w-full text-center text-sm text-gray-500 dark:text-gray-400 hover:text-gray-700 dark:hover:text-gray-200 transition-colors"
          >
            Cancelar e voltar
          </button>
        </div>
      </div>
    </div>
  );
}