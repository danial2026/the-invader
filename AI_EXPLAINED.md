# Behavioral Architecture & Tactical Logic

> *"Every purpose is established by counsel: and with good advice make war."* — **Proverbs 20:18**

This document details the logic systems and behavioral patterns implemented within **The Invader**. The system is comprised of two core modules: the **Automated Pilot** and the **Entity Tactical Logic**.

---

## The Automated Pilot

The Automated Pilot is a deterministic, rule-based system designed for optimal engagement efficiency. It operates via a hierarchical decision matrix:

1.  **Threat Mitigation**: The system continuously scans for incoming kinetic threats. If a projectile is detected within the terminal danger radius (30 pixels), the system initiates an immediate evasive burn.
2.  **Target Prioritization**: When clear of immediate threats, the pilot identifies the optimal target based on proximity and engagement probability.
3.  **Variable Engagement**: The pilot maintains a stabilized firing rate, capped at 3 concurrent projectiles to ensure system throughput and target clarity.

---

## Entity Behavioral Profiles

Tactical variety is achieved through a dynamic mapping system that assigns entities one of seven specialized behavioral profiles based on their biological traits.

### Logic Mapping
The system analyzes raw biographical data to assign one of the following tactical frameworks:
- **Combatant/Hostile** ➔ Aggressive
- **Cautious/Anxious** ➔ Defensive
- **Unpredictable** ➔ Chaotic
- **Tracker** ➔ Hunter
- **Observer** ➔ Stalker
- **Vanguard** ➔ Berserker
- **Precisionist** ➔ Sniper

### Tactical Frameworks

| Profile | Operational Speed | Tactical Objective |
|---------|-------------------|--------------------|
| **Aggressive** | High (5.0) | Direct intercept; utilizes periodic vertical compression to overwhelm defenses. |
| **Defensive** | Low (2.0) | Prioritizes evasion; maintains safe standoff distance from the player. |
| **Chaotic** | Variable (4.0) | Employs randomized directional shifts and vertical oscillations to complicate target tracking. |
| **Hunter** | Predictive (4.0) | Utilizes predictive lead-tracking to intercept player movement. |
| **Stalker** | Low (1.0) | Employs slow approach vectors; stabilizes completely before firing for maximum precision. |
| **Berserker** | Oscillatory (6.0) | High-speed sine-wave movement; executes vertical drops upon perimeter impact. |
| **Sniper** | Static (1.0) | Positions at tactical anchor points; maintains high-accuracy suppression from a distance. |

---

## Core System Mechanics

### Proximity Deconfliction
To ensure operational efficiency, entities utilize an anti-collision algorithm. If the distance between two units drops below 30 pixels, the system applies a corrective repulsion force to maintain formation integrity.

### Optimized Logic Processing
By utilizing a deterministic rule-based architecture rather than a neural-network-based model, the behavioral engine maintains high performance with minimal computational overhead, ensuring a consistent 60Hz update cycle.
