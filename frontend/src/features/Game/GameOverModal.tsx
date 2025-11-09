

import { useGameStore } from "../../stores/useGameStore";

export function GameModal() {
  const { gameState } = useGameStore();
  if (!gameState || gameState.status.code == "online") return;

  let modal : React.ReactNode;
  switch (gameState.status.result) {
    case "win":
      modal = <GameWinModal />
      break
    case "draw":
      modal = <GameDrawModal />
      break
    case "loss":
      modal = <GameLossModal />
      break
  }

  return (
    <div className="absolute bg-black/50 w-full h-full">
      {modal}
    </div>
    
  )
}

function GameWinModal() {
  return (
    <>
      <div className="absolute inset-0 w-1/2 h-1/2">
        YOU WON YAYAYAYAYAYAYAYAY
      
      </div>
    </>
  )
}



function GameLossModal() {
  return (
    <>
      <div className="absolute inset-0 w-1/2 h-1/2">
        YOU LOST HAHAHAHAHAHAHAHAHA
      
      </div>
    
    </>
  )
}

function GameDrawModal() {
  return (
    <>
      <div className="absolute inset-0 w-1/2 h-1/2">
        WOMP WOMP NOBODY WON
      
      </div>
    
    
    </>
  )
}