# LinkedIn Automation (PoC)

An advanced educational proof-of-concept for browser automation using Go and Rod. This project demonstrates stealth automation techniques, structural DOM scraping, and robust session management.

**⚠️ Disclaimer:** This project is for technical evaluation and educational purposes only. Automated interaction with LinkedIn violates their Terms of Service. Do not use on primary accounts.

**Project google doc link**

    https://docs.google.com/document/d/1w3_9AgOkgalhf7JvJ92QkIyHa-e1CBxFc937CBiA6OU/edit?usp=sharing

## Key Features

- **Stealth Mode**: Implements human-like behavior including Bézier curve mouse movements, random jitter, and variable typing speeds.
- **Smart Session Management**: Persists cookies locally to avoid repeated logins. Detects valid sessions and skips authentication automatically.
- **Robust Scraper**: Uses structural selectors (DOM-agnostic) rather than brittle class names, ensuring reliability against UI updates.
- **Hybrid Outreach**: Intelligently handles different profile layouts (e.g., "+ Add" buttons, "Connect" buttons, or hidden "More" dropdowns).
- **Safety Locks**: Includes cooldowns and limit checks to mimic natural usage.

## Requirements

- **Go** 1.20+
- **Google Chrome** installed locally
- **Git**

## Setup & Run

1.  **Clone the repository**

        ```bash
        git clone <repo-url>
        cd linkedin-automation
        ```

        2.**Configure Environment Create a .env file in the project root:**

        LINKEDIN_USERNAME=your_email@example.com
        LINKEDIN_PASSWORD=your_secure_password

        3.**Install Dependencies**

        go mod tidy

        4.**Run the Bot**

        go run cmd/bot/main.go

            **Configuration**

        Credentials: Loaded via .env.

        Targets: Search keywords and location are currently defined in cmd/bot/main.go. You can update them there to change the target role/city.

        Limits: The bot is currently capped at 5 profiles per run for safety.


        **Project Structure**

        Path,Description

    cmd/bot/main.go,Application entry point and main control flow.
    internal/auth,"Smart login logic, session persistence, and security checkpoint detection."
    internal/bot,Browser initialization and stealth injection.
    internal/search,Search engine logic and resilient result parsing.
    internal/outreach,"Logic for finding ""Connect"" buttons, handling modals, and sending notes."
    internal/stealth,"Human simulation helpers (Mouse movement, Typing)."
    internal/models,"Shared data structures (e.g., Lead)."
    internal/config,Configuration management.

    **Troubleshooting (Windows)**
    Chrome Zombie Processes
    If the bot fails to attach or launch, old Chrome instances might be blocking the port. Open PowerShell and run:

    taskkill /F /IM chrome.exe /T

    **Reset Session**
    To force a fresh login (e.g., to switch accounts or fix corrupted cookies), delete the temporary profile folder:

    C:\temp\rod-profile
