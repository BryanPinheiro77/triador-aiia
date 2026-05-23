"use client";

import { FormEvent, useCallback, useEffect, useState } from "react";

type Analysis = {
  id: number;
  candidate_name: string;
  skills: string[];
  years_experience: number;
  fit_score: number;
  summary: string;
};

const API_URL = process.env.NEXT_PUBLIC_API_URL ?? "http://localhost:8080";

export default function Home() {
  const [resume, setResume] = useState("");
  const [jobDescription, setJobDescription] = useState("");
  const [result, setResult] = useState<Analysis | null>(null);
  const [animatedScore, setAnimatedScore] = useState(0);
  const [history, setHistory] = useState<Analysis[]>([]);
  const [loading, setLoading] = useState(false);
  const [historyLoading, setHistoryLoading] = useState(false);
  const [error, setError] = useState("");


  const loadHistory = useCallback(async () => {
    setHistoryLoading(true);

    try {
      const response = await fetch(`${API_URL}/analyses`);

      if (!response.ok) {
        throw new Error("Não foi possível carregar o histórico.");
      }

      const data: Analysis[] = await response.json();
      setHistory(data);
    } catch (error) {
      setError(
        error instanceof Error
          ? error.message
          : "Erro inesperado ao carregar histórico."
      );
    } finally {
      setHistoryLoading(false);
    }
  }, []);

  useEffect(() => {
    // eslint-disable-next-line react-hooks/set-state-in-effect
    void loadHistory();
  }, [loadHistory]);

  useEffect(() => {
  if (!result) {
    // eslint-disable-next-line react-hooks/set-state-in-effect
    return;
  }

  const timeout = setTimeout(() => {
    setAnimatedScore(result.fit_score);
  }, 100);

  return () => clearTimeout(timeout);
}, [result]);

  async function handleSubmit(event: FormEvent<HTMLFormElement>) {
    event.preventDefault();

    setError("");
    setResult(null);
    setAnimatedScore(0);

    if (!resume.trim() || !jobDescription.trim()) {
      setError("Preencha o currículo e a descrição da vaga.");
      return;
    }

    setLoading(true);

    try {
      const response = await fetch(`${API_URL}/analyses`, {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
        },
        body: JSON.stringify({
          resume,
          job_description: jobDescription,
        }),
      });

      const data = await response.json();

      if (!response.ok) {
        throw new Error(data.error ?? "Erro ao analisar currículo.");
      }

      setResult(data);
      setResume("");
      setJobDescription("");

      await loadHistory();
    } catch (error) {
      setError(
        error instanceof Error
          ? error.message
          : "Erro inesperado ao analisar currículo."
      );
    } finally {
      setLoading(false);
    }
  }

  return (
    <main className="min-h-screen bg-slate-950 px-6 py-10 text-slate-100">
      <div className="mx-auto max-w-6xl space-y-8">
        <header className="space-y-3">
          <p className="text-sm font-semibold uppercase tracking-[0.3em] text-blue-400">
            Triador
          </p>

          <div className="space-y-4">
            <h1 className="max-w-3xl text-4xl font-bold tracking-tight md:text-5xl">
              Análise de aderência entre currículo e vaga
            </h1>

            <p className="max-w-3xl text-base leading-7 text-slate-300">
              Cole o texto de um currículo e de uma vaga para gerar uma análise
              estruturada via IA, com score, skills, experiência estimada e
              resumo justificando a nota.
            </p>
          </div>
        </header>

        <section className="grid gap-6 lg:grid-cols-[1.2fr_0.8fr]">
          <form
            onSubmit={handleSubmit}
            className="space-y-5 rounded-2xl border border-slate-800 bg-slate-900/70 p-6 shadow-xl"
          >
            <div className="space-y-2">
              <label
                htmlFor="resume"
                className="text-sm font-medium text-slate-200"
              >
                Currículo
              </label>

              <textarea
                id="resume"
                value={resume}
                onChange={(event) => setResume(event.target.value)}
                placeholder="Cole aqui o texto do currículo..."
                className="min-h-48 w-full resize-y rounded-xl border border-slate-700 bg-slate-950 p-4 text-sm text-slate-100 outline-none transition focus:border-blue-500"
              />
            </div>

            <div className="space-y-2">
              <label
                htmlFor="jobDescription"
                className="text-sm font-medium text-slate-200"
              >
                Descrição da vaga
              </label>

              <textarea
                id="jobDescription"
                value={jobDescription}
                onChange={(event) => setJobDescription(event.target.value)}
                placeholder="Cole aqui a descrição da vaga..."
                className="min-h-40 w-full resize-y rounded-xl border border-slate-700 bg-slate-950 p-4 text-sm text-slate-100 outline-none transition focus:border-blue-500"
              />
            </div>

            {error && (
              <div className="rounded-xl border border-red-900 bg-red-950/60 p-4 text-sm text-red-200">
                {error}
              </div>
            )}

            <button
              type="submit"
              disabled={loading}
              className="w-full rounded-xl bg-blue-500 px-5 py-3 text-sm font-semibold text-white transition hover:bg-blue-400 disabled:cursor-not-allowed disabled:bg-slate-700 disabled:text-slate-400"
            >
              {loading ? (
                <span className="flex items-center justify-center gap-2">
                  <span className="h-4 w-4 animate-spin rounded-full border-2 border-white/40 border-t-white" />
                  Analisando com IA...
                </span>
              ) : (
               "Analisar currículo"
            )}
            </button>
            {loading && (
              <p className="text-center text-sm text-slate-400">
                A IA está comparando o currículo com a vaga. Isso pode levar alguns segundos.
               </p>
            )}
          </form>

          <aside className="rounded-2xl border border-slate-800 bg-slate-900/70 p-6 shadow-xl">
            <h2 className="text-xl font-semibold">Resultado</h2>

            {!result ? (
              <p className="mt-4 text-sm leading-6 text-slate-400">
                O resultado da análise aparecerá aqui após o envio do formulário.
              </p>
            ) : (
              <div className="mt-5 space-y-5">
                <div>
                  <p className="text-sm text-slate-400">Candidato</p>
                  <p className="text-lg font-semibold">
                    {result.candidate_name}
                  </p>
                </div>

<div>
  <p className="text-sm text-slate-400">
    Score de aderência
  </p>

  <p className="text-4xl font-bold text-blue-400">
    {result.fit_score}
    <span className="text-lg text-slate-400">
      /100
    </span>
  </p>

  <div className="mt-3 h-3 w-full rounded-full bg-slate-800">
    <div
      className="h-3 rounded-full bg-blue-500 transition-all duration-700 ease-out"
      style={{
        width: `${animatedScore}%`,
      }}
    />
  </div>
</div>

                <div>
                  <p className="text-sm text-slate-400">Skills</p>
                  <div className="mt-2 flex flex-wrap gap-2">
                    {result.skills.map((skill) => (
                      <span
                        key={skill}
                        className="rounded-full bg-blue-500/10 px-3 py-1 text-xs font-medium text-blue-300"
                      >
                        {skill}
                      </span>
                    ))}
                  </div>
                </div>

                <div>
                  <p className="text-sm text-slate-400">Resumo</p>
                  <p className="mt-1 text-sm leading-6 text-slate-200">
                    {result.summary}
                  </p>
                </div>
              </div>
            )}
          </aside>
        </section>

        <section className="rounded-2xl border border-slate-800 bg-slate-900/70 p-6 shadow-xl">
          <div className="flex items-center justify-between gap-4">
            <div>
              <h2 className="text-xl font-semibold">Histórico</h2>
              <p className="mt-1 text-sm text-slate-400">
                Últimas análises salvas no banco.
              </p>
            </div>

            <button
              type="button"
              onClick={loadHistory}
              disabled={historyLoading}
              className="rounded-xl border border-slate-700 px-4 py-2 text-sm font-medium text-slate-200 transition hover:border-blue-500 hover:text-blue-300 disabled:cursor-not-allowed disabled:text-slate-500"
            >
              {historyLoading ? "Carregando..." : "Atualizar"}
            </button>
          </div>

          <div className="mt-6 grid gap-4 md:grid-cols-2">
            {history.length === 0 && !historyLoading ? (
              <div className="col-span-full rounded-xl border border-dashed border-slate-700 bg-slate-950/60 p-6 text-center">
                <p className="text-sm font-medium text-slate-200">
                  Nenhuma análise encontrada
                </p>
                <p className="mt-2 text-sm text-slate-400">
                  Envie um currículo e uma descrição de vaga para criar a primeira análise.
                </p>
              </div>
            ) : (
              history.map((analysis) => (
                <article
                  key={analysis.id}
                  className="rounded-xl border border-slate-800 bg-slate-950 p-4"
                >
                  <div className="flex items-start justify-between gap-3">
                    <div>
                      <h3 className="font-semibold">
                        {analysis.candidate_name}
                      </h3>
                      <p className="mt-1 text-sm text-slate-400">
                        {analysis.years_experience} ano(s) de experiência
                      </p>
                    </div>

                    <span className="rounded-full bg-blue-500/10 px-3 py-1 text-sm font-semibold text-blue-300">
                      {analysis.fit_score}/100
                    </span>
                  </div>

                  <p className="mt-3 text-sm leading-6 text-slate-300">
                    {analysis.summary}
                  </p>

                  <div className="mt-3 flex flex-wrap gap-2">
                    {analysis.skills.slice(0, 6).map((skill) => (
                      <span
                        key={`${analysis.id}-${skill}`}
                        className="rounded-full bg-slate-800 px-2.5 py-1 text-xs text-slate-300"
                      >
                        {skill}
                      </span>
                    ))}
                  </div>
                </article>
              ))
            )}
          </div>
        </section>
      </div>
    </main>
  );
}