
export interface Participant {
    name: string
}
export interface Span {
    upper:  number
    lower:  number
}
export interface NodeReference {
    id: number
    kind: number
}

export interface MatchResult {
    winningNode: number
    winningSlot: number 
    winningPlayer: number 
    losingParticipant: number 
    winningParticipant: number 
    dropNode: number
}
export interface GameState {
    result: MatchResult

} 
export interface Node {
    id: number
    nodeType: number
    span: Span
    gridSpan: Span
    drop:   number
    participant: number
    challengerUpOne: boolean
    state: GameState
}


export interface Bracket {

    root:   Node
    rootNodeId:  number
    nodes:  Node[]
    span:  Span
    gridSpan:   Span
    level: number
    tier: number
    drop: number
    preferences: Preferences
    participants: number[] 
}

export interface Tournament {
    id: number
    typ: number
    name: string
    bracket:  Bracket
    participants: Participant[]
    tournamentState: string

}

export interface Selection {
    node: Game
}

export as namespace brackets;
