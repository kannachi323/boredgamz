import { AppLayout } from './Layout';
import { Outlet } from 'react-router-dom';

function App() {

  return (
    <main>
      <AppLayout>
        <Outlet />
      </AppLayout>
    </main>
  );
}


export default App;