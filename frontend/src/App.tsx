import { AppLayout } from './Layout';
import { Outlet } from 'react-router-dom';

import { GameProvider } from './contexts/GameProvider';


function App() {

  return (
    <main>
      <GameProvider>
        <AppLayout>
          <Outlet />
        </AppLayout>
      </GameProvider>
    </main>
  );
}


export default App;