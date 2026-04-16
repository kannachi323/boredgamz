import { AppLayout } from './Layout';
import { Outlet, useLocation } from 'react-router-dom';
import { SiteLockedPage } from './components/SiteLockedPage';
import { SITE_LOCKDOWN_ENABLED } from './config/siteAccess';

function App() {
  const location = useLocation();
  const isHomePage = location.pathname === '/';
  const showLockedPage = SITE_LOCKDOWN_ENABLED && !isHomePage;

  return (
    <main>
      <AppLayout>
        {showLockedPage ? <SiteLockedPage /> : <Outlet />}
      </AppLayout>
    </main>
  );
}


export default App;
