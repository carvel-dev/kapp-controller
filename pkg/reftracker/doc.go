// Package reftracker contains structs used for tracking secret and configmap referenced by the app.
// Tracking these references allows us to trigger an app reconcile when the resources are updated.
package reftracker
