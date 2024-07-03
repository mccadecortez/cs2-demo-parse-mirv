import constants from './mirv/constants.mjs';
import { MirvJS } from './mirv/mirv.mjs';
import { events } from './mirv/ws-events.mjs';
{
	MirvJS.init({
		host: constants.websocketListenHost,
		port: constants.websocketListenPort,
		path: constants.websocketListenPath
	});
	MirvJS.wsEnable = true;
	// main logic defined here, it runs on every tick
	mirv.onClientFrameStageNotify = (e) => {
		// FRAME_START - called on host_frame (1 per tick).
		if (e.curStage == 0 && e.isBefore) {
			MirvJS.connect();
			if (MirvJS.ws !== null) {
				// Flush any messages that are lingering:
				MirvJS.ws.flush();
			}
			MirvJS.tick++;
		}
		// FRAME_RENDER_START - this is not called when demo is paused (can be multiple per tick).
		if (e.curStage === 5 && e.isBefore) {
			if (MirvJS.ws !== null && mirv.isPlayingDemo()) {
				try {
					MirvJS.ws.send(
						JSON.stringify({
							type: events.demoTick,
							data: { server_tick: mirv.getDemoTick() }
						})
					);

					// we could flush and then wait for a reply here to set a view instantly, but don't understimate network round-trip time!
				} catch (err) {
					mirv.warning(
						'frame-demo-tick: Error while sending message:' + String(err) + '\n'
					);
				}
			}
		}
		// FRAME_RENDER_END - this is not called when demo is paused (can be multiple per tick).
		// note: double check if it's 6
		if (e.curStage === 5 && e.isBefore) {
			if (MirvJS.ws !== null) MirvJS.ws.flush();
		}
	};
}
