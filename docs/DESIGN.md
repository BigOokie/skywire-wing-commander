# High level design: skywire-telegram-notify-bot
This document provides the high-level design outline for the skywire-telegram-notify-bot application (written in Go).
All key elements of the Bot application are covered within this document (or will be once completed - this is a work in progress).

# Purpose
The Bot is intended to run on a Skywire Manager Node and provide its owner with near real-time notifications about significant events that occur within the Skywire Nodes that are managed by the Skywire Manager (collectivley referred to as the SkyMiner).

# Key Features
The following key features are to be supported by the Bot applications. Those denoted with an `!` are not yet implemented and are roadmap concepts currently - however consideration has been given to these (at some level) within the overall design of the Bot application.

## Automated notifications
* **Heartbeat** notifications. Used to notify the owner of the Bot that the Bot application is still alive and running. Can be turned on or off (default).
* Skywire Node **inbound** and **outbound** **connection events**. Monitoring of the Skywire Nodes can be turned on and off (default). When enabled, the Bot will send an alert to the associated Telegram chat session (managed by Telegram API ChatID). A summary of the changes will be shown.

## Ad-hoc (manual) interaction
* **Status** notifications can be requested from the bot at any time. The response is the same as the automated Heartbeat (when enabled). This allow the owner of the Bot to ask if it is still running and responsive to commands.
* **Help**. Provides a list of commands supported by the Bot.
* **About**. Provides details of the author and any contributors, the Bot application version, and links to the GitHub repository.

# Design
## Overview
The Bot application is designed with a single main command and control processing loop which is responsible for management and co-ordination of all other parts of the Bot application.

Once the main processing loop is initiated, it monitors for Telegram chat messages. Any message that is not a pre-defined command known to the Bot is discarded and ignored.

Any known command (designated with the `/`) will be dispatched to a command processor. The command processor is responsible for initiating the behaviour associated with each command.

Generally speaking, ad-hock commands are managed interactivly within the scope of the command processors process. Automated commands are managed by background processes which are co-ordinated and controlled by the main process.

## Skywire Node Monitoring
This section describes the high level approach to monitoring Skywire Nodes from the Manager Node.

The early releases of the Bot (`v0.0.1-alpha` thru `v0.0.4-alpha`) utilises OS level file monitoring events to observe changes made to the `clients.json` file. While this did provide some degree of alerting for connection changes - it was not sufficient, and later versions utilise the inbuilt API provided within the Skywire Manager and Nodes.

The following discusses the planned approach for versions of the Bot App starting with `v0.0.5-alpha`. Some of this remains as conceptual design and may not yet be implemented.

### Manager and Node API Usage
The following APIs provided by the Skywire Manager and Node applications are utilised.
* **Manager**:
    * `getToken`
    * `login`
    * `getAll`
    * `getNode`
* **Node**:
    * `getInfo`
    * `getApps`

### API Workflows
This section outlines the usage patterns for the Manager and Node APIs.

* Call `login` providing the Manager password. This will authenticate the Bot App with the Manager Node (and API)
* Call `getToken`. The obtains an authentication token for the Bot App from the Manager node. The authentication token needs to be provided in aother requests.
* Call `getAll` to obtain the list of Nodes connected to the Manager. This is expected to be a list of zero, one or more Nodes.
* Itterate the list of connected nodes (provided by `getAll`). For each Node:
    * Call `getNode` to ontain specific detail about the current Node. This is needed to make API calls on the specific Node.
    * For the current Node:
        * Call `getInfo` to get general information back from the Node.
        * Call `getApps` to get a list of Skywire pplications that are currently run (and being managed) by the Node.

Information obtained from the above needs to be cached by the Bot. This represents the current state of the SkyMiner that is being monitored. This information can be used to determine if changes in connected state change, and what clients and applications are running. Basically the Bot will have access to the information that is availabe within the Skywire Manager Web GUI.

The Bot App will periodically poll the Manager and Node APIs  described above. When specific conditions change, the Bot will send a notification to the connected Telegram chat. Specific events that will trigger a notification from the Bot app include:
* The list of Nodes (from `getAll`) changes. The total number of connected Nodes will be reported in the Telegram chat.
* The state of a specific Node changes (from `getInfo`). This will signify that a connection has been made to that Node (either inbound or outbound)

## Future Ideas
The following section outlines some conceptual ideas that should be possible.

* Allow the user to enquire on the status of a specific Node. It should be possible for the user to ask for the status of a Node. The Bot should prompt the user to select one of the connected Nodes from a list. When the user selects a specific Node, an enquiry for that Node is then performed.
