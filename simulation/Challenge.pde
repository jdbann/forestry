abstract class Challenge {
  World world;
  EventStream stream;
  boolean complete;
  String title;
  int completedAt;

  PGraphics overlay;

  Challenge() {
    stream = new EventStream(this);
    int seed = int(random(256*256));
    world = new World(seed, stream);
    complete = false;

    world.addTrees(treeCount);
    world.addLogs(logCount);
  }

  void tick() {
    world.tick();
  }

  PGraphics overlay() {
    if (overlay == null) {
      overlay = drawOverlay();
    }
    return overlay;
  }

  abstract void listen(Event e);

  abstract void openHelp();

  PGraphics drawOverlay() {
    overlay = createGraphics(width, height, P3D);
    overlay.beginDraw();
    overlay.textSize(32);
    overlay.fill(255);
    overlay.text(title, 32, 64);
    if (complete) {
      overlay.text("Complete!", 32, 96, 400, height);
    } else {
      overlay.textSize(16);
    }
    overlay.endDraw();
    return overlay;
  }

  void markComplete() {
    if (complete == false) {
      complete = true;
      overlay = null;
      completedAt = frameCount;
    }
  }
}
