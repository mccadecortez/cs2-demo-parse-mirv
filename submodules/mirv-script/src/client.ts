import { SimpleWebSocket } from 'simple-websockets';
import { EntityObject } from './mirv/utils.mjs';
import { events } from './mirv/ws-events.mjs';
import { MirvMessage } from './server.js';

type ConnectionOptions = {
	host: string;
	port: number;
	path?: string;
	user: number | string;
};

export class MirvClient {
	private ws: SimpleWebSocket;

	/**
	 * @param path - path is optional, if not provided it will default to mirv.
	 */
	constructor({ host, port, path, user }: ConnectionOptions) {
		this.ws = new SimpleWebSocket(`ws://${host}:${port}/${path || 'mirv'}?user=${user}`);
		this.ws.on(events.warning, (message) => {
			console.log('warning:', message);
		});
	}

	send(message: MirvMessage) {
		this.ws.send(message.type, message.data);
	}
	/** List available types */
	listTypes(callback: (types: string[]) => void) {
		this.ws.once(events.listTypes, callback);
		this.send({ type: 'listTypes' });
	}
	/** Shorthand for quit */
	sendQuit() {
		this.send({ type: 'quit' });
	}
	/** Execute ingame command */
	sendExec(command: string) {
		this.send({ type: 'exec', data: command });
	}
	/** Enable transimission of game events */
	enableGameEvents(callback: (gameEvents: mirv.GameEvent) => void) {
		this.ws.on(events.onGameEvent, callback);
		this.send({ type: 'setGameEvents', data: true });
	}
	/** Disable transimission of game events */
	disableGameEvents() {
		this.ws.removeAllListeners(events.onGameEvent);
		this.send({ type: 'setGameEvents', data: false });
	}
	/** Enable transimission of cViewRenderSetupView events */
	enableCViewRenderSetupView(callback: (view: mirv.OnCViewRenderSetupViewArgs) => void) {
		this.ws.on(events.onCViewRenderSetupView, callback);
		this.send({ type: 'setCViewRenderSetupView', data: true });
	}
	/** Disable transimission of cViewRenderSetupView events */
	disableCViewRenderSetupView() {
		this.ws.removeAllListeners(events.onCViewRenderSetupView);
		this.send({ type: 'setCViewRenderSetupView', data: false });
	}
	/** Get last cached render view  */
	getLastView(callback: (view: mirv.OnCViewRenderSetupViewArgs['lastView']) => void) {
		this.ws.once(events.getLastView, callback);
		this.send({ type: 'getLastView' });
	}
	/** Set render view */
	setView(view: mirv.OnCViewRenderSetupViewSet) {
		this.send({ type: 'setView', data: view });
	}
	/** Reset render view */
	resetView() {
		this.send({ type: 'setView' });
	}
	/** Enable transimission of entity events */
	enableEntityEvents(
		onAddEntity: (entity: EntityObject) => void,
		onRemoveEntity: (entity: EntityObject) => void
	) {
		this.ws.on(events.onAddEntity, onAddEntity);
		this.ws.on(events.onRemoveEntity, onRemoveEntity);
		this.send({ type: 'setEntityEvents', data: true });
	}
	/** Disable transimission of entity events */
	disableEntityEvents() {
		this.ws.removeAllListeners(events.onAddEntity);
		this.ws.removeAllListeners(events.onRemoveEntity);
		this.send({ type: 'setEntityEvents', data: false });
	}
	/** List all non null entities */
	listEntities(callback: (entities: EntityObject[]) => void) {
		this.ws.once(events.listEntities, callback);
		this.send({ type: 'listEntities' });
	}
	/** List all player entities */
	listPlayerEntities(callback: (entities: EntityObject[]) => void) {
		this.ws.once(events.listPlayerEntities, callback);
		this.send({ type: 'listPlayerEntities' });
	}
	/** Load a module
	 * @param fullPath - full path to the module `.mjs` file. It should exist on the same machine as the CS2/HLAE
	 */
	loadModule(fullPath: string) {
		this.send({ type: 'loadModule', data: fullPath });
	}
}

// test
const client = new MirvClient({ host: 'localhost', port: 31337, user: 1 });
