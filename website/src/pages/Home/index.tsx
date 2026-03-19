
import { MainSection } from './MainSection';
import { FeaturedSection } from './FeaturedSection';
import { CommunitySection } from './CommunitySection';
import { FooterV1 } from '../../features/Footer/FooterV1';

export default function App() {

  return (
    <div className="min-h-screen bg-gray-900 font-sans antialiased">
    
      <main>
        <MainSection />
        <FeaturedSection />
        <CommunitySection />
      </main>

      <FooterV1 />
    </div>
  );
}