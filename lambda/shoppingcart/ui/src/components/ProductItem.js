import React from 'react'
import PropTypes from 'prop-types'
import { Table } from 'semantic-ui-react'

const ProductItem = ({ product, onAddToCartClicked }) => (
    <Table.Row key={'product/' + product.id}>
    <Table.Cell>{product.name}</Table.Cell>
    <Table.Cell textAlign='right'>{product.price}</Table.Cell>
    <Table.Cell textAlign='right'>{product.quantity}</Table.Cell>
    <Table.Cell textAlign='right'>
    <button
        onClick={onAddToCartClicked}
        disabled={product.quantity > 0 ? '' : 'disabled'}>
        {product.quantity > 0 ? 'Add to cart' : 'Sold Out'}
    </button>
    </Table.Cell>
    </Table.Row>
)

ProductItem.propTypes = {
  product: PropTypes.shape({
    name: PropTypes.string.isRequired,
    price: PropTypes.number.isRequired,
    quantity: PropTypes.number.isRequired
  }).isRequired,
  onAddToCartClicked: PropTypes.func.isRequired
}

export default ProductItem
