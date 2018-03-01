# Zelduh

A monochrome 4-color game!

## Goals

- [ ] State machines
    - [x] State machine for UI
        - [x] Start screen
        - [x] Game screen
        - [x] Pause screen
        - [x] Game over screen
    - [ ] State machine for actual game (maybe multiple)
- [ ] Player character movement
    - [x] Move with arrow keys
    - [x] Pressing and holding keys for continuous movement
    - [x] Edge of window boundaries
    - [ ] Get rid of delay on win.Repeated()
- [ ] NPC movement
    - [ ] AI
        * Look at "steering patterns"
        * What behaviors do I want characters to have and how can this be achieved realistically
        * Keep it flexible
        * Iterative process
        * Units should react in situations that seem realistic in an environment they seem to be interactin with
        * Don't box in what enemies can do too early
    - [ ] If player character is in line of sight, enemies charge/shoot
- [ ] Projectiles
    - [ ] Player sword is a projectile that shoots one tile 
    - [ ] Enemies can shoot projectiles that go x tiles
- [ ] Animated sprites
    - [ ] Pixel art
    - [ ] Sprite sheet
- [ ] Collision detection
    - [ ] Enemies cause damage
    - [ ] Obstacles
    - [ ] Block projectiles with shield
    - [ ] Pickup items
- [ ] Interact with environment
    - [ ] Open things
    - [ ] Push things
- [ ] Basic stats
    - [ ] Health
    - [ ] Charge sword
- [ ] Start fully equiped
- [ ] High score goal
    - [ ] Score points by killing enemies, collecting items
- [ ] Attack
    - [ ] Basic sword slash
    - [ ] Charged magic sword projectile
