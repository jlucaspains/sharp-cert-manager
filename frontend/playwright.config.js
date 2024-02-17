import { devices } from '@playwright/test';

/** @type {import('@playwright/test').PlaywrightTestConfig} */
const config = {
	testDir: 'tests',
	testMatch: /(.+\.)?(test|spec)\.[jt]s/,
	use: {
		actionTimeout: 0,
		baseURL: process.env.BASEURL || 'http://localhost:5173/',
	
		trace: 'on-first-retry',
		video: 'off',
		screenshot: 'only-on-failure',
	  },
	
	  projects: [
		{
		  name: 'chromium',
		  use: {
			...devices['Desktop Chrome'],
		  },
		},
		{
		  name: 'webkit',
		  use: {
			...devices['Desktop Safari'],
		  },
		},
		{
		  name: 'firefox',
		  use: {
			...devices['DeskDesktop Firefox'],
		  },
		},
		{
		  name: 'Mobile Chrome',
		  use: {
			...devices['Pixel 5'],
		  },
		},
		{
		  name: 'Mobile Safari',
		  use: {
			...devices['iPhone 12'],
		  },
		},
	  ],
};

export default config;
