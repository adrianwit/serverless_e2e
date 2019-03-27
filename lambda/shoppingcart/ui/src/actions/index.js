import shop from '../api/shop'
import * as types from '../constants/ActionTypes'



const receiveProducts = (resp) => {
    if(resp.Status === "ok") {
        return {
            type: types.RECEIVE_PRODUCTS_SUCCESS,
            products: resp.Data
        }
    }
    return {
        type: types.RECEIVE_PRODUCTS_FAILURE,
        error: resp.Error
    }
}


export const getAllProducts = () => dispatch => {
    shop.getProducts(resp => {
        dispatch(receiveProducts(resp))
    })
}




const addToCartUnsafe = productId => ({
    type: types.ADD_TO_CART,
    productId
})



export const addToCart = productId => (dispatch, getState) => {
    if (getState().products.byId[productId].quantity > 0) {

        dispatch(addToCartUnsafe(productId))



    }
}

export const checkout = products => (dispatch, getState) => {
    const {cart} = getState()

    dispatch({
        type: types.CHECKOUT_REQUEST
    })


    shop.buyProducts(products, () => {
        dispatch({
            type: types.CHECKOUT_SUCCESS,
            cart
        })
        // Replace the line above with line below to rollback on failure:
        // dispatch({ type: types.CHECKOUT_FAILURE, cart })
    })
}
