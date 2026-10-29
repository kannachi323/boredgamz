import { SITE_LOCKDOWN_MESSAGE } from "../config/siteAccess";

export function HomeOverlayNotice() {
  return (
    <div className="pointer-events-none fixed inset-0 z-40 flex items-center justify-center bg-black/55 px-6">
      <div className="pointer-events-auto w-full max-w-2xl rounded-3xl border border-gray-700 bg-gray-900/95 p-8 text-center shadow-2xl backdrop-blur-sm sm:p-10">
        <p className="text-sm font-semibold uppercase tracking-[0.3em] text-green-400">
          Temporary Notice
        </p>
        <h2 className="mt-4 text-3xl font-black tracking-tight text-white sm:text-4xl">
          Quick update
        </h2>
        <p className="mt-5 text-lg leading-8 text-gray-300">
          {SITE_LOCKDOWN_MESSAGE}
        </p>
      </div>
    </div>
  );
}
