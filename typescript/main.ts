import { parse } from 'https://deno.land/std@0.207.0/flags/mod.ts';

// Import the CLI module
if (import.meta.main) {
  // Parse command line arguments to determine which module to run
  const args = parse(Deno.args, {
    string: ['mode'],
    default: { mode: 'cli' },
  });

  // Run the appropriate module based on the mode
  switch (args.mode) {
    case 'cli':
      // Import and run the CLI module
      import('./src/cli/pm_chat.ts').then(module => {
        console.log('Starting Personal Agent CLI...');
      }).catch(error => {
        console.error('Failed to start CLI:', error);
      });
      break;
    default:
      console.error(`Unknown mode: ${args.mode}`);
      Deno.exit(1);
  }
}
