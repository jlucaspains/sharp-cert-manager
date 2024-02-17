import { expect, test } from '@playwright/test';

test('index page has expected h1', async ({ page }) => {
	await page.goto('/');
	await expect(page.getByText('sharp-cert-manager')).toBeVisible();
});

test('list of sites load', async ({ page }) => {
	await page.goto('/');
	await expect(page.getByTestId('result-item')).toHaveCount(6)
});

test('open a specific item', async ({ page }) => {
	await page.goto('/');

	await page.getByRole('button', { name: 'lpains.net Issuer: GTS CA 1D4' }).click();
	await expect(page.getByRole('cell', { name: 'lpains.net' }).first()).toBeVisible();
	await expect(page.getByRole('cell', { name: 'SHA256-RSA' })).toBeVisible();
});
