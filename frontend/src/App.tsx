import { AppLayout } from './Layout';
import { Outlet } from 'react-router-dom';

import { GameProvider } from './contexts/GameProvider';
import { AuthProvider } from './contexts/AuthProvider';


function App() {

  return (
    <main>
      <AuthProvider>
        <GameProvider>
          <AppLayout>
            <Outlet />
          </AppLayout>
        </GameProvider>
      </AuthProvider>
    </main>
  );
}


export default App;