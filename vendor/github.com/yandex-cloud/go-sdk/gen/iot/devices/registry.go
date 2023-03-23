// Code generated by sdkgen. DO NOT EDIT.

// nolint
package devices

import (
	"context"

	"google.golang.org/grpc"

	devices "github.com/yandex-cloud/go-genproto/yandex/cloud/iot/devices/v1"
	"github.com/yandex-cloud/go-genproto/yandex/cloud/operation"
)

//revive:disable

// RegistryServiceClient is a devices.RegistryServiceClient with
// lazy GRPC connection initialization.
type RegistryServiceClient struct {
	getConn func(ctx context.Context) (*grpc.ClientConn, error)
}

// AddCertificate implements devices.RegistryServiceClient
func (c *RegistryServiceClient) AddCertificate(ctx context.Context, in *devices.AddRegistryCertificateRequest, opts ...grpc.CallOption) (*operation.Operation, error) {
	conn, err := c.getConn(ctx)
	if err != nil {
		return nil, err
	}
	return devices.NewRegistryServiceClient(conn).AddCertificate(ctx, in, opts...)
}

// AddDataStreamExport implements devices.RegistryServiceClient
func (c *RegistryServiceClient) AddDataStreamExport(ctx context.Context, in *devices.AddDataStreamExportRequest, opts ...grpc.CallOption) (*operation.Operation, error) {
	conn, err := c.getConn(ctx)
	if err != nil {
		return nil, err
	}
	return devices.NewRegistryServiceClient(conn).AddDataStreamExport(ctx, in, opts...)
}

// AddPassword implements devices.RegistryServiceClient
func (c *RegistryServiceClient) AddPassword(ctx context.Context, in *devices.AddRegistryPasswordRequest, opts ...grpc.CallOption) (*operation.Operation, error) {
	conn, err := c.getConn(ctx)
	if err != nil {
		return nil, err
	}
	return devices.NewRegistryServiceClient(conn).AddPassword(ctx, in, opts...)
}

// Create implements devices.RegistryServiceClient
func (c *RegistryServiceClient) Create(ctx context.Context, in *devices.CreateRegistryRequest, opts ...grpc.CallOption) (*operation.Operation, error) {
	conn, err := c.getConn(ctx)
	if err != nil {
		return nil, err
	}
	return devices.NewRegistryServiceClient(conn).Create(ctx, in, opts...)
}

// Delete implements devices.RegistryServiceClient
func (c *RegistryServiceClient) Delete(ctx context.Context, in *devices.DeleteRegistryRequest, opts ...grpc.CallOption) (*operation.Operation, error) {
	conn, err := c.getConn(ctx)
	if err != nil {
		return nil, err
	}
	return devices.NewRegistryServiceClient(conn).Delete(ctx, in, opts...)
}

// DeleteCertificate implements devices.RegistryServiceClient
func (c *RegistryServiceClient) DeleteCertificate(ctx context.Context, in *devices.DeleteRegistryCertificateRequest, opts ...grpc.CallOption) (*operation.Operation, error) {
	conn, err := c.getConn(ctx)
	if err != nil {
		return nil, err
	}
	return devices.NewRegistryServiceClient(conn).DeleteCertificate(ctx, in, opts...)
}

// DeleteDataStreamExport implements devices.RegistryServiceClient
func (c *RegistryServiceClient) DeleteDataStreamExport(ctx context.Context, in *devices.DeleteDataStreamExportRequest, opts ...grpc.CallOption) (*operation.Operation, error) {
	conn, err := c.getConn(ctx)
	if err != nil {
		return nil, err
	}
	return devices.NewRegistryServiceClient(conn).DeleteDataStreamExport(ctx, in, opts...)
}

// DeletePassword implements devices.RegistryServiceClient
func (c *RegistryServiceClient) DeletePassword(ctx context.Context, in *devices.DeleteRegistryPasswordRequest, opts ...grpc.CallOption) (*operation.Operation, error) {
	conn, err := c.getConn(ctx)
	if err != nil {
		return nil, err
	}
	return devices.NewRegistryServiceClient(conn).DeletePassword(ctx, in, opts...)
}

// Get implements devices.RegistryServiceClient
func (c *RegistryServiceClient) Get(ctx context.Context, in *devices.GetRegistryRequest, opts ...grpc.CallOption) (*devices.Registry, error) {
	conn, err := c.getConn(ctx)
	if err != nil {
		return nil, err
	}
	return devices.NewRegistryServiceClient(conn).Get(ctx, in, opts...)
}

// GetByName implements devices.RegistryServiceClient
func (c *RegistryServiceClient) GetByName(ctx context.Context, in *devices.GetByNameRegistryRequest, opts ...grpc.CallOption) (*devices.Registry, error) {
	conn, err := c.getConn(ctx)
	if err != nil {
		return nil, err
	}
	return devices.NewRegistryServiceClient(conn).GetByName(ctx, in, opts...)
}

// List implements devices.RegistryServiceClient
func (c *RegistryServiceClient) List(ctx context.Context, in *devices.ListRegistriesRequest, opts ...grpc.CallOption) (*devices.ListRegistriesResponse, error) {
	conn, err := c.getConn(ctx)
	if err != nil {
		return nil, err
	}
	return devices.NewRegistryServiceClient(conn).List(ctx, in, opts...)
}

type RegistryIterator struct {
	ctx  context.Context
	opts []grpc.CallOption

	err           error
	started       bool
	requestedSize int64
	pageSize      int64

	client  *RegistryServiceClient
	request *devices.ListRegistriesRequest

	items []*devices.Registry
}

func (c *RegistryServiceClient) RegistryIterator(ctx context.Context, req *devices.ListRegistriesRequest, opts ...grpc.CallOption) *RegistryIterator {
	var pageSize int64
	const defaultPageSize = 1000
	pageSize = req.PageSize
	if pageSize == 0 {
		pageSize = defaultPageSize
	}
	return &RegistryIterator{
		ctx:      ctx,
		opts:     opts,
		client:   c,
		request:  req,
		pageSize: pageSize,
	}
}

func (it *RegistryIterator) Next() bool {
	if it.err != nil {
		return false
	}
	if len(it.items) > 1 {
		it.items[0] = nil
		it.items = it.items[1:]
		return true
	}
	it.items = nil // consume last item, if any

	if it.started && it.request.PageToken == "" {
		return false
	}
	it.started = true

	if it.requestedSize == 0 || it.requestedSize > it.pageSize {
		it.request.PageSize = it.pageSize
	} else {
		it.request.PageSize = it.requestedSize
	}

	response, err := it.client.List(it.ctx, it.request, it.opts...)
	it.err = err
	if err != nil {
		return false
	}

	it.items = response.Registries
	it.request.PageToken = response.NextPageToken
	return len(it.items) > 0
}

func (it *RegistryIterator) Take(size int64) ([]*devices.Registry, error) {
	if it.err != nil {
		return nil, it.err
	}

	if size == 0 {
		size = 1 << 32 // something insanely large
	}
	it.requestedSize = size
	defer func() {
		// reset iterator for future calls.
		it.requestedSize = 0
	}()

	var result []*devices.Registry

	for it.requestedSize > 0 && it.Next() {
		it.requestedSize--
		result = append(result, it.Value())
	}

	if it.err != nil {
		return nil, it.err
	}

	return result, nil
}

func (it *RegistryIterator) TakeAll() ([]*devices.Registry, error) {
	return it.Take(0)
}

func (it *RegistryIterator) Value() *devices.Registry {
	if len(it.items) == 0 {
		panic("calling Value on empty iterator")
	}
	return it.items[0]
}

func (it *RegistryIterator) Error() error {
	return it.err
}

// ListCertificates implements devices.RegistryServiceClient
func (c *RegistryServiceClient) ListCertificates(ctx context.Context, in *devices.ListRegistryCertificatesRequest, opts ...grpc.CallOption) (*devices.ListRegistryCertificatesResponse, error) {
	conn, err := c.getConn(ctx)
	if err != nil {
		return nil, err
	}
	return devices.NewRegistryServiceClient(conn).ListCertificates(ctx, in, opts...)
}

type RegistryCertificatesIterator struct {
	ctx  context.Context
	opts []grpc.CallOption

	err           error
	started       bool
	requestedSize int64
	pageSize      int64

	client  *RegistryServiceClient
	request *devices.ListRegistryCertificatesRequest

	items []*devices.RegistryCertificate
}

func (c *RegistryServiceClient) RegistryCertificatesIterator(ctx context.Context, req *devices.ListRegistryCertificatesRequest, opts ...grpc.CallOption) *RegistryCertificatesIterator {
	var pageSize int64
	const defaultPageSize = 1000

	if pageSize == 0 {
		pageSize = defaultPageSize
	}
	return &RegistryCertificatesIterator{
		ctx:      ctx,
		opts:     opts,
		client:   c,
		request:  req,
		pageSize: pageSize,
	}
}

func (it *RegistryCertificatesIterator) Next() bool {
	if it.err != nil {
		return false
	}
	if len(it.items) > 1 {
		it.items[0] = nil
		it.items = it.items[1:]
		return true
	}
	it.items = nil // consume last item, if any

	if it.started {
		return false
	}
	it.started = true

	response, err := it.client.ListCertificates(it.ctx, it.request, it.opts...)
	it.err = err
	if err != nil {
		return false
	}

	it.items = response.Certificates
	return len(it.items) > 0
}

func (it *RegistryCertificatesIterator) Take(size int64) ([]*devices.RegistryCertificate, error) {
	if it.err != nil {
		return nil, it.err
	}

	if size == 0 {
		size = 1 << 32 // something insanely large
	}
	it.requestedSize = size
	defer func() {
		// reset iterator for future calls.
		it.requestedSize = 0
	}()

	var result []*devices.RegistryCertificate

	for it.requestedSize > 0 && it.Next() {
		it.requestedSize--
		result = append(result, it.Value())
	}

	if it.err != nil {
		return nil, it.err
	}

	return result, nil
}

func (it *RegistryCertificatesIterator) TakeAll() ([]*devices.RegistryCertificate, error) {
	return it.Take(0)
}

func (it *RegistryCertificatesIterator) Value() *devices.RegistryCertificate {
	if len(it.items) == 0 {
		panic("calling Value on empty iterator")
	}
	return it.items[0]
}

func (it *RegistryCertificatesIterator) Error() error {
	return it.err
}

// ListDataStreamExports implements devices.RegistryServiceClient
func (c *RegistryServiceClient) ListDataStreamExports(ctx context.Context, in *devices.ListDataStreamExportsRequest, opts ...grpc.CallOption) (*devices.ListDataStreamExportsResponse, error) {
	conn, err := c.getConn(ctx)
	if err != nil {
		return nil, err
	}
	return devices.NewRegistryServiceClient(conn).ListDataStreamExports(ctx, in, opts...)
}

type RegistryDataStreamExportsIterator struct {
	ctx  context.Context
	opts []grpc.CallOption

	err           error
	started       bool
	requestedSize int64
	pageSize      int64

	client  *RegistryServiceClient
	request *devices.ListDataStreamExportsRequest

	items []*devices.DataStreamExport
}

func (c *RegistryServiceClient) RegistryDataStreamExportsIterator(ctx context.Context, req *devices.ListDataStreamExportsRequest, opts ...grpc.CallOption) *RegistryDataStreamExportsIterator {
	var pageSize int64
	const defaultPageSize = 1000

	if pageSize == 0 {
		pageSize = defaultPageSize
	}
	return &RegistryDataStreamExportsIterator{
		ctx:      ctx,
		opts:     opts,
		client:   c,
		request:  req,
		pageSize: pageSize,
	}
}

func (it *RegistryDataStreamExportsIterator) Next() bool {
	if it.err != nil {
		return false
	}
	if len(it.items) > 1 {
		it.items[0] = nil
		it.items = it.items[1:]
		return true
	}
	it.items = nil // consume last item, if any

	if it.started {
		return false
	}
	it.started = true

	response, err := it.client.ListDataStreamExports(it.ctx, it.request, it.opts...)
	it.err = err
	if err != nil {
		return false
	}

	it.items = response.DataStreamExports
	return len(it.items) > 0
}

func (it *RegistryDataStreamExportsIterator) Take(size int64) ([]*devices.DataStreamExport, error) {
	if it.err != nil {
		return nil, it.err
	}

	if size == 0 {
		size = 1 << 32 // something insanely large
	}
	it.requestedSize = size
	defer func() {
		// reset iterator for future calls.
		it.requestedSize = 0
	}()

	var result []*devices.DataStreamExport

	for it.requestedSize > 0 && it.Next() {
		it.requestedSize--
		result = append(result, it.Value())
	}

	if it.err != nil {
		return nil, it.err
	}

	return result, nil
}

func (it *RegistryDataStreamExportsIterator) TakeAll() ([]*devices.DataStreamExport, error) {
	return it.Take(0)
}

func (it *RegistryDataStreamExportsIterator) Value() *devices.DataStreamExport {
	if len(it.items) == 0 {
		panic("calling Value on empty iterator")
	}
	return it.items[0]
}

func (it *RegistryDataStreamExportsIterator) Error() error {
	return it.err
}

// ListDeviceTopicAliases implements devices.RegistryServiceClient
func (c *RegistryServiceClient) ListDeviceTopicAliases(ctx context.Context, in *devices.ListDeviceTopicAliasesRequest, opts ...grpc.CallOption) (*devices.ListDeviceTopicAliasesResponse, error) {
	conn, err := c.getConn(ctx)
	if err != nil {
		return nil, err
	}
	return devices.NewRegistryServiceClient(conn).ListDeviceTopicAliases(ctx, in, opts...)
}

type RegistryDeviceTopicAliasesIterator struct {
	ctx  context.Context
	opts []grpc.CallOption

	err           error
	started       bool
	requestedSize int64
	pageSize      int64

	client  *RegistryServiceClient
	request *devices.ListDeviceTopicAliasesRequest

	items []*devices.DeviceAlias
}

func (c *RegistryServiceClient) RegistryDeviceTopicAliasesIterator(ctx context.Context, req *devices.ListDeviceTopicAliasesRequest, opts ...grpc.CallOption) *RegistryDeviceTopicAliasesIterator {
	var pageSize int64
	const defaultPageSize = 1000
	pageSize = req.PageSize
	if pageSize == 0 {
		pageSize = defaultPageSize
	}
	return &RegistryDeviceTopicAliasesIterator{
		ctx:      ctx,
		opts:     opts,
		client:   c,
		request:  req,
		pageSize: pageSize,
	}
}

func (it *RegistryDeviceTopicAliasesIterator) Next() bool {
	if it.err != nil {
		return false
	}
	if len(it.items) > 1 {
		it.items[0] = nil
		it.items = it.items[1:]
		return true
	}
	it.items = nil // consume last item, if any

	if it.started && it.request.PageToken == "" {
		return false
	}
	it.started = true

	if it.requestedSize == 0 || it.requestedSize > it.pageSize {
		it.request.PageSize = it.pageSize
	} else {
		it.request.PageSize = it.requestedSize
	}

	response, err := it.client.ListDeviceTopicAliases(it.ctx, it.request, it.opts...)
	it.err = err
	if err != nil {
		return false
	}

	it.items = response.Aliases
	it.request.PageToken = response.NextPageToken
	return len(it.items) > 0
}

func (it *RegistryDeviceTopicAliasesIterator) Take(size int64) ([]*devices.DeviceAlias, error) {
	if it.err != nil {
		return nil, it.err
	}

	if size == 0 {
		size = 1 << 32 // something insanely large
	}
	it.requestedSize = size
	defer func() {
		// reset iterator for future calls.
		it.requestedSize = 0
	}()

	var result []*devices.DeviceAlias

	for it.requestedSize > 0 && it.Next() {
		it.requestedSize--
		result = append(result, it.Value())
	}

	if it.err != nil {
		return nil, it.err
	}

	return result, nil
}

func (it *RegistryDeviceTopicAliasesIterator) TakeAll() ([]*devices.DeviceAlias, error) {
	return it.Take(0)
}

func (it *RegistryDeviceTopicAliasesIterator) Value() *devices.DeviceAlias {
	if len(it.items) == 0 {
		panic("calling Value on empty iterator")
	}
	return it.items[0]
}

func (it *RegistryDeviceTopicAliasesIterator) Error() error {
	return it.err
}

// ListOperations implements devices.RegistryServiceClient
func (c *RegistryServiceClient) ListOperations(ctx context.Context, in *devices.ListRegistryOperationsRequest, opts ...grpc.CallOption) (*devices.ListRegistryOperationsResponse, error) {
	conn, err := c.getConn(ctx)
	if err != nil {
		return nil, err
	}
	return devices.NewRegistryServiceClient(conn).ListOperations(ctx, in, opts...)
}

type RegistryOperationsIterator struct {
	ctx  context.Context
	opts []grpc.CallOption

	err           error
	started       bool
	requestedSize int64
	pageSize      int64

	client  *RegistryServiceClient
	request *devices.ListRegistryOperationsRequest

	items []*operation.Operation
}

func (c *RegistryServiceClient) RegistryOperationsIterator(ctx context.Context, req *devices.ListRegistryOperationsRequest, opts ...grpc.CallOption) *RegistryOperationsIterator {
	var pageSize int64
	const defaultPageSize = 1000
	pageSize = req.PageSize
	if pageSize == 0 {
		pageSize = defaultPageSize
	}
	return &RegistryOperationsIterator{
		ctx:      ctx,
		opts:     opts,
		client:   c,
		request:  req,
		pageSize: pageSize,
	}
}

func (it *RegistryOperationsIterator) Next() bool {
	if it.err != nil {
		return false
	}
	if len(it.items) > 1 {
		it.items[0] = nil
		it.items = it.items[1:]
		return true
	}
	it.items = nil // consume last item, if any

	if it.started && it.request.PageToken == "" {
		return false
	}
	it.started = true

	if it.requestedSize == 0 || it.requestedSize > it.pageSize {
		it.request.PageSize = it.pageSize
	} else {
		it.request.PageSize = it.requestedSize
	}

	response, err := it.client.ListOperations(it.ctx, it.request, it.opts...)
	it.err = err
	if err != nil {
		return false
	}

	it.items = response.Operations
	it.request.PageToken = response.NextPageToken
	return len(it.items) > 0
}

func (it *RegistryOperationsIterator) Take(size int64) ([]*operation.Operation, error) {
	if it.err != nil {
		return nil, it.err
	}

	if size == 0 {
		size = 1 << 32 // something insanely large
	}
	it.requestedSize = size
	defer func() {
		// reset iterator for future calls.
		it.requestedSize = 0
	}()

	var result []*operation.Operation

	for it.requestedSize > 0 && it.Next() {
		it.requestedSize--
		result = append(result, it.Value())
	}

	if it.err != nil {
		return nil, it.err
	}

	return result, nil
}

func (it *RegistryOperationsIterator) TakeAll() ([]*operation.Operation, error) {
	return it.Take(0)
}

func (it *RegistryOperationsIterator) Value() *operation.Operation {
	if len(it.items) == 0 {
		panic("calling Value on empty iterator")
	}
	return it.items[0]
}

func (it *RegistryOperationsIterator) Error() error {
	return it.err
}

// ListPasswords implements devices.RegistryServiceClient
func (c *RegistryServiceClient) ListPasswords(ctx context.Context, in *devices.ListRegistryPasswordsRequest, opts ...grpc.CallOption) (*devices.ListRegistryPasswordsResponse, error) {
	conn, err := c.getConn(ctx)
	if err != nil {
		return nil, err
	}
	return devices.NewRegistryServiceClient(conn).ListPasswords(ctx, in, opts...)
}

type RegistryPasswordsIterator struct {
	ctx  context.Context
	opts []grpc.CallOption

	err           error
	started       bool
	requestedSize int64
	pageSize      int64

	client  *RegistryServiceClient
	request *devices.ListRegistryPasswordsRequest

	items []*devices.RegistryPassword
}

func (c *RegistryServiceClient) RegistryPasswordsIterator(ctx context.Context, req *devices.ListRegistryPasswordsRequest, opts ...grpc.CallOption) *RegistryPasswordsIterator {
	var pageSize int64
	const defaultPageSize = 1000

	if pageSize == 0 {
		pageSize = defaultPageSize
	}
	return &RegistryPasswordsIterator{
		ctx:      ctx,
		opts:     opts,
		client:   c,
		request:  req,
		pageSize: pageSize,
	}
}

func (it *RegistryPasswordsIterator) Next() bool {
	if it.err != nil {
		return false
	}
	if len(it.items) > 1 {
		it.items[0] = nil
		it.items = it.items[1:]
		return true
	}
	it.items = nil // consume last item, if any

	if it.started {
		return false
	}
	it.started = true

	response, err := it.client.ListPasswords(it.ctx, it.request, it.opts...)
	it.err = err
	if err != nil {
		return false
	}

	it.items = response.Passwords
	return len(it.items) > 0
}

func (it *RegistryPasswordsIterator) Take(size int64) ([]*devices.RegistryPassword, error) {
	if it.err != nil {
		return nil, it.err
	}

	if size == 0 {
		size = 1 << 32 // something insanely large
	}
	it.requestedSize = size
	defer func() {
		// reset iterator for future calls.
		it.requestedSize = 0
	}()

	var result []*devices.RegistryPassword

	for it.requestedSize > 0 && it.Next() {
		it.requestedSize--
		result = append(result, it.Value())
	}

	if it.err != nil {
		return nil, it.err
	}

	return result, nil
}

func (it *RegistryPasswordsIterator) TakeAll() ([]*devices.RegistryPassword, error) {
	return it.Take(0)
}

func (it *RegistryPasswordsIterator) Value() *devices.RegistryPassword {
	if len(it.items) == 0 {
		panic("calling Value on empty iterator")
	}
	return it.items[0]
}

func (it *RegistryPasswordsIterator) Error() error {
	return it.err
}

// Update implements devices.RegistryServiceClient
func (c *RegistryServiceClient) Update(ctx context.Context, in *devices.UpdateRegistryRequest, opts ...grpc.CallOption) (*operation.Operation, error) {
	conn, err := c.getConn(ctx)
	if err != nil {
		return nil, err
	}
	return devices.NewRegistryServiceClient(conn).Update(ctx, in, opts...)
}