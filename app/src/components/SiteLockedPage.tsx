import { SITE_LOCKDOWN_MESSAGE } from "../config/siteAccess";

export function SiteLockedPage() {
  return (
    <section className="min-h-full bg-gray-950 px-6 py-12 text-white">
      <div className="mx-auto flex min-h-[calc(92vh-6rem)] max-w-3xl items-center justify-center">
        <div className="w-full rounded-3xl border border-gray-800 bg-gray-900/90 p-8 text-center shadow-2xl sm:p-12">
          <p className="text-sm font-semibold uppercase tracking-[0.3em] text-green-400">
            Temporary hold
          </p>
          <h1 className="mt-4 text-4xl font-black tracking-tight sm:text-5xl">
            This page is offline for now
          </h1>
          <p className="mx-auto mt-6 max-w-2xl text-lg leading-8 text-gray-300">
            {SITE_LOCKDOWN_MESSAGE}
          </p>
          <div className="mt-8 flex justify-center">
            <a
              href="/"
              className="rounded-xl bg-green-400 px-6 py-3 font-semibold text-gray-950 transition hover:bg-green-300"
            >
              Back to home
            </a>
          </div>
        </div>
      </div>
    </section>
  );
}
