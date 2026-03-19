import { FaDiscord, FaGithub } from "react-icons/fa";
import { MdOutlineEmail } from "react-icons/md"

export function CommunitySection() {
  return (
    <div className="text-center flex flex-col justify-evenly items-center gap-10 m-10">
      <h2 className="text-3xl md:text-5xl font-bold text-text tracking-tight flex flex-col gap-2">
        <span>Join our community</span>
        <span className="text-text-muted text-lg md:text-2xl font-medium">
          Meet players, share feedback, and stay in the loop.
        </span>
      </h2>


      <span className="flex flex-row items-center gap-5">
        <a
          href="/community"
          className="px-8 py-3 rounded-full bg-zinc-100 hover:bg-zinc-100/90 text-black font-semibold text-sm"
        >
          Join Community
        </a>

      
        <div className="flex items-center gap-4 px-5 py-2.5 rounded-full border border-white/10 bg-white/5 backdrop-blur-sm">
          <span className="text-sm text-text-muted">Connect with us</span>
          <div className="w-px h-4 bg-white/10" />
          <div className="flex items-center gap-3">
            <a href="https://discord.gg/boredgamz" target="_blank" rel="noopener noreferrer" className="text-text-muted hover:text-zinc-200 transition-colors" aria-label="Discord">
              <FaDiscord className="text-3xl" />
            </a>
            <a href="https://github.com/your-repo/discussions" target="_blank" rel="noopener noreferrer" className="text-text-muted hover:text-zinc-200 transition-colors" aria-label="GitHub Discussions">
              <FaGithub className="text-2xl" />
            </a>
            <a href="mailto:hello@boredgamz.com" className="text-text-muted hover:text-zinc-200 transition-colors" aria-label="Email">
              <MdOutlineEmail className="text-3xl" />
            </a>
          </div>
        </div>


      </span>

      
    </div>
  );
}